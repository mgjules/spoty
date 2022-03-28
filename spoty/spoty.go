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
	URL    string     `json:"url"`
	Height int        `json:"height"`
	Width  int        `json:"width"`
	RGBA   color.RGBA `json:"rgba,omitempty"`
	Hex    string     `json:"hex,omitempty"`
	Error  string     `json:"error,omitempty"`
}

// Spoty represents the spoty service.
type Spoty struct {
	client *spotify.Client

	auth  spotify.Authenticator
	state string

	cache  *cache.Cache
	health *health.Checks
}

// New creates a new spoty service.
func New(cfg *config.Config, cache *cache.Cache, health *health.Checks) (*Spoty, error) {
	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return nil, errors.New("missing clientID or clientSecret")
	}

	auth := spotify.NewAuthenticator(
		fmt.Sprintf("http://%s:%d/api/callback", cfg.Host, cfg.Port),
		spotify.ScopeUserReadCurrentlyPlaying,
		spotify.ScopeUserReadPlaybackState,
	)

	auth.SetAuthInfo(cfg.ClientID, cfg.ClientSecret)

	state, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("new uuid: %w", err)
	}

	spoty := Spoty{
		auth:   auth,
		state:  state.String(),
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
func (s *Spoty) TrackCurrentlyPlaying() (*spotify.FullTrack, error) {
	const cacheCurrentTrackKey = "current_track"

	cachedTrack, found := s.cache.Get(cacheCurrentTrackKey)
	if found {
		if cachedTrack, ok := cachedTrack.(*spotify.FullTrack); ok {
			return cachedTrack, nil
		}
	}

	if !s.IsPlaying() {
		return nil, errors.New("no track currently playing")
	}

	playing, err := s.client.PlayerCurrentlyPlaying()
	if err != nil {
		return nil, err
	}

	s.cache.SetWithTTL(cacheCurrentTrackKey, playing.Item, 0, _defaultTTL)

	return playing.Item, nil
}

// TrackImages returns the track images from a track.
func (s *Spoty) TrackImages(track *spotify.FullTrack) ([]Image, error) {
	if track == nil {
		return nil, errors.New("invalid track")
	}

	cacheTrackImagesKey := "track_" + strcase.ToCamel(string(track.ID)) + "_images"

	cachedImages, found := s.cache.Get(cacheTrackImagesKey)
	if found {
		if cachedImages, ok := cachedImages.([]Image); ok {
			return cachedImages, nil
		}
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
				wg.Done()
			}()

			resp, err := httpClient.Get(albumImage.URL)
			if err != nil {
				img.Error = "could not retrieve album image"

				return
			}
			defer resp.Body.Close() //nolint: errcheck

			processedImg, _, err := image.Decode(resp.Body)
			if err != nil {
				img.Error = "could not process album image"

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
