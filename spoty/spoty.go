package spoty

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"sync"
	"time"

	"github.com/JulesMike/spoty/cache"
	"github.com/JulesMike/spoty/config"
	"github.com/JulesMike/spoty/health"
	"github.com/JulesMike/spoty/logger"
	"github.com/JulesMike/spoty/tracer"
	"github.com/cenkalti/dominantcolor"
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
	"github.com/zmb3/spotify"
	"go.uber.org/fx"
)

const _defaultTTL = 5 * time.Second

// Module exported for initialising a new Spoty service.
var Module = fx.Options(
	fx.Provide(New),
)

// Image represents an image with its dominant color.
type Image struct {
	URL      string     `json:"url"`
	Height   int        `json:"height"`
	Width    int        `json:"width"`
	RGBA     color.RGBA `json:"rgba,omitempty"`
	Hex      string     `json:"hex,omitempty"`
	Error    string     `json:"error,omitempty"`
	RawError error      `json:"-"`
}

// Spoty represents the spoty service.
type Spoty struct {
	client *spotify.Client

	auth  spotify.Authenticator
	state string

	logger *logger.Logger
	tracer *tracer.Tracer
	cache  *cache.Cache
	health *health.Checks
}

// New creates a new spoty service.
func New(
	cfg *config.Config,
	logger *logger.Logger,
	tracer *tracer.Tracer,
	cache *cache.Cache,
	health *health.Checks,
) (*Spoty, error) {
	if cfg.SpotifyClientID == "" || cfg.SpotifyClientSecret == "" {
		return nil, errors.New("missing clientID or clientSecret")
	}

	auth := spotify.NewAuthenticator(
		fmt.Sprintf("http://%s:%d/api/callback", cfg.HttpServerHost, cfg.HttpServerPort),
		spotify.ScopeUserReadCurrentlyPlaying,
		spotify.ScopeUserReadPlaybackState,
	)

	auth.SetAuthInfo(cfg.SpotifyClientID, cfg.SpotifyClientSecret)

	state, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("new uuid: %w", err)
	}

	spoty := Spoty{
		auth:   auth,
		state:  state.String(),
		logger: logger,
		tracer: tracer,
		cache:  cache,
		health: health,
	}

	spoty.health.RegisterChecks(spoty.Check())

	return &spoty, nil
}

// IsAuth returns true if the client was created.
func (s *Spoty) IsAuth() bool {
	return s.client != nil
}

// IsPlaying returns true if the client is currently playing.
func (s *Spoty) IsPlaying() bool {
	state, err := s.client.PlayerState()
	if err != nil {
		return false
	}

	return state.Playing
}

// AuthURL returns the spotify auth url.
func (s *Spoty) AuthURL() string {
	return s.auth.AuthURL(s.state)
}

// SetupNewClient sets up a new spotify client.
func (s *Spoty) SetupNewClient(r *http.Request) error {
	tok, err := s.auth.Token(s.state, r)
	if err != nil {
		return err
	}

	client := s.auth.NewClient(tok)
	client.AutoRetry = true
	s.client = &client

	return nil
}

// TrackCurrentlyPlaying returns the currently playing track.
func (s *Spoty) TrackCurrentlyPlaying(ctx context.Context) (*spotify.FullTrack, error) {
	ctx, span := s.tracer.Start(ctx, "TrackCurrentlyPlaying")
	defer span.End()

	const cacheCurrentTrackKey = "current_track"

	cachedTrack, found := s.cache.Get(cacheCurrentTrackKey)
	if found {
		if cachedTrack, ok := cachedTrack.(*spotify.FullTrack); ok {
			s.logger.Ctx(ctx).Debugw("found cached track", "track", cachedTrack)

			return cachedTrack, nil
		}

		s.logger.Ctx(ctx).Debugw("failed to parse cached track. retrieving fresh one...", "track", cachedTrack)
	}

	if !s.IsPlaying() {
		s.logger.ErrorwContext(ctx, "no track currently playing")

		return nil, errors.New("no track currently playing")
	}

	playing, err := s.client.PlayerCurrentlyPlaying()
	if err != nil {
		s.logger.ErrorwContext(ctx, "failed to retrieve currently playing track", "error", err.Error())

		return nil, err
	}

	s.cache.SetWithTTL(cacheCurrentTrackKey, playing.Item, 0, _defaultTTL)

	return playing.Item, nil
}

// TrackImages returns the track images from a track.
func (s *Spoty) TrackImages(ctx context.Context, track *spotify.FullTrack) ([]Image, error) {
	ctx, span := s.tracer.Start(ctx, "TrackImages")
	defer span.End()

	if track == nil {
		return nil, errors.New("invalid track")
	}

	cacheTrackImagesKey := "track_" + strcase.ToCamel(string(track.ID)) + "_images"

	cachedImages, found := s.cache.Get(cacheTrackImagesKey)
	if found {
		if cachedImages, ok := cachedImages.([]Image); ok {
			s.logger.Ctx(ctx).Debugw("found cached images", "images", cachedImages)

			return cachedImages, nil
		}

		s.logger.Ctx(ctx).Debugw("failed to parse cached images. retrieving fresh ones...", "track", cachedImages)
	}

	httpClient := &http.Client{
		Timeout: _defaultTTL,
	}

	var wg sync.WaitGroup

	var images []Image
	for i := range track.Album.Images {
		albumImage := &track.Album.Images[i]

		wg.Add(1)
		go func(albumImage *spotify.Image) {
			img := Image{
				URL:    albumImage.URL,
				Height: albumImage.Height,
				Width:  albumImage.Width,
			}

			defer func() {
				images = append(images, img)

				if img.Error != "" {
					s.logger.WarnwContext(ctx, img.Error, "error", img.RawError.Error(), "image", img)
				}

				wg.Done()
			}()

			resp, err := httpClient.Get(img.URL)
			if err != nil {
				img.Error = "could not retrieve album image"
				img.RawError = err

				return
			}
			defer resp.Body.Close() //nolint: errcheck

			processedImg, _, err := image.Decode(resp.Body)
			if err != nil {
				img.Error = "could not process album image"
				img.RawError = err

				return
			}

			img.RGBA = dominantcolor.Find(processedImg)
			img.Hex = dominantcolor.Hex(img.RGBA)
		}(albumImage)
	}

	wg.Wait()

	s.cache.SetWithTTL(cacheTrackImagesKey, images, 0, _defaultTTL)

	return images, nil
}

// Check checks if the spoty service is authenticated.
func (s *Spoty) Check() health.Check {
	//nolint:revive
	return health.Check{
		Name:          "spoty",
		RefreshPeriod: 10 * time.Second,
		InitialDelay:  10 * time.Second,
		Timeout:       5 * time.Second,
		Check: func(ctx context.Context) error {
			if s.IsAuth() {
				return nil
			}

			return errors.New("spoty not authenticated")
		},
	}
}
