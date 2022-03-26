package spoty

import (
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
	"github.com/cenkalti/dominantcolor"
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
	"github.com/zmb3/spotify"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(New),
)

type Image struct {
	URL    string     `json:"url"`
	Height int        `json:"height"`
	Width  int        `json:"width"`
	RGBA   color.RGBA `json:"rgba,omitempty"`
	Hex    string     `json:"hex,omitempty"`
	Error  string     `json:"error,omitempty"`
}

type Spoty struct {
	client *spotify.Client

	auth  spotify.Authenticator
	state string

	cache *cache.Cache
}

func New(cfg *config.Config, cache *cache.Cache) (*Spoty, error) {
	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return nil, errors.New("missing clientID or clientSecret")
	}

	if cache == nil {
		return nil, errors.New("invalid cache")
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

	return &Spoty{
		auth:  auth,
		state: state.String(),
		cache: cache,
	}, nil
}

func (s *Spoty) IsAuth() bool {
	return s.client != nil
}

func (s *Spoty) IsPlaying() bool {
	state, err := s.client.PlayerState()
	if err != nil {
		return false
	}

	return state.Playing
}

func (s *Spoty) AuthURL() string {
	return s.auth.AuthURL(s.state)
}

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

func (s *Spoty) TrackCurrentlyPlaying() (*spotify.FullTrack, error) {
	const cacheCurrentTrackKey = "current_track"

	cachedTrack, found := s.cache.Get(cacheCurrentTrackKey)
	if found {
		return cachedTrack.(*spotify.FullTrack), nil
	}

	if !s.IsPlaying() {
		return nil, errors.New("no track currently playing")
	}

	playing, err := s.client.PlayerCurrentlyPlaying()
	if err != nil {
		return nil, err
	}

	s.cache.SetWithTTL(cacheCurrentTrackKey, playing.Item, 0, 5*time.Second)

	return playing.Item, nil
}

func (s *Spoty) TrackImages(track *spotify.FullTrack) ([]Image, error) {
	if track == nil {
		return nil, errors.New("invalid track")
	}

	cacheTrackImagesKey := "track_" + strcase.ToCamel(string(track.ID)) + "_images"

	cachedImages, found := s.cache.Get(cacheTrackImagesKey)
	if found {
		return cachedImages.([]Image), nil
	}

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	var wg sync.WaitGroup

	var images []Image
	for _, albumImage := range track.Album.Images {
		wg.Add(1)
		go func(albumImage spotify.Image) {
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
			defer resp.Body.Close()

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

	s.cache.SetWithTTL(cacheTrackImagesKey, images, 0, 5*time.Second)

	return images, nil
}
