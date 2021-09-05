package ldtk

import (
	"bytes"
	"embed"
	"encoding/json"
	"image"
	"log"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
	Px   [2]int `json:"px"`
	Src  [2]int `json:"src"`
	Type int    `json:"t"`
	Flip int    `json:"f"`
}

type IntGridTile struct {
	Coord int `json:"coordId"`
	Value int `json:"v"`
}

type Entity struct {
	ID     string   `json:"__identifier"`
	Width  int      `json:"width"`
	Height int      `json:"height"`
	Px     [2]int   `json:"px"`
	Fields []*Field `json:"fieldInstances"`
}

type Neighbour struct {
	UID int    `json:"levelUid"`
	Dir string `json:"dir"`
}

type Field struct {
	ID    string      `json:"__identifier"`
	Type  string      `json:"__type"`
	Value interface{} `json:"__value"`
}

type Layer struct {
	ID             string         `json:"__identifier"`
	Type           string         `json:"__type"`
	Entities       []*Entity      `json:"entityInstances"`
	Size           int            `json:"__gridSize"`
	Width          int            `json:"__cWid"`
	Height         int            `json:"__cHei"`
	AutoLayerTiles []*Tile        `json:"autoLayerTiles"`
	GridTiles      []*Tile        `json:"gridTiles"`
	Tileset        string         `json:"__tilesetRelPath"`
	IntGrid        []*IntGridTile `json:"intGrid"`
	Tiles          map[int]*ebiten.Image
}

type Level struct {
	ID         string       `json:"identifier"`
	UID        int          `json:"uid"`
	Fields     []*Field     `json:"fieldInstances"`
	Layers     []*Layer     `json:"layerInstances"`
	BgColor    string       `json:"__bgColor"`
	Width      int          `json:"pxWid"`
	Height     int          `json:"pxHei"`
	Neighbours []*Neighbour `json:"__neighbours"`
}

type Ldtk struct {
	Levels  []*Level `json:"levels"`
	Active  *Level
	Options *Options
}

type Options struct {
	// a pointer to the folder containing assets
	EmbedFolder *embed.FS
	// whether LDtk is using aseprite files or png files
	Aseprite bool
	// the root path to the map assets (ltdk and png)
	Root string
	// LDtk int grid values for which we set up collisions
	CollidesWith map[int]bool
	// this is called every time we encounter a tile we should set up for
	OnCollisionAdd func(x float64, y float64, w float64, h float64)
}

func Load(path string, opt *Options) (*Ldtk, error) {
	data, err := opt.EmbedFolder.ReadFile(opt.Root + "/" + path)
	if err != nil {
		log.Fatal(err)
	}
	var ldtkMap Ldtk
	ldtkMap.Options = opt
	err = json.Unmarshal(data, &ldtkMap)
	if err != nil {
		return nil, err
	}
	for _, level := range ldtkMap.Levels {
		level.PrepareLevel(*opt)
	}
	ldtkMap.Active = ldtkMap.Levels[0]
	ldtkMap.Active.MakeCollisions(*opt)
	return &ldtkMap, nil
}

func (f *Field) AsArray() []interface{} {
	return f.Value.([]interface{})
}

func (f *Field) AsString() string {
	return f.Value.(string)
}

func (f *Field) AsInt() int {
	return int(f.AsFloat64())
}

func (f *Field) AsFloat64() float64 {
	return f.Value.(float64)
}

func (f *Field) AsBool() bool {
	return f.Value.(bool)
}

func (f *Field) AsMap() map[string]interface{} {
	return f.Value.(map[string]interface{})
}

// PrepareLevel populates the layer.Tiles prop if required
func (level *Level) PrepareLevel(opt Options) {
	for _, layer := range level.Layers {
		if len(layer.AutoLayerTiles) > 0 {
			layer.MakeTiles(layer.AutoLayerTiles, opt)
		} else if len(layer.GridTiles) > 0 {
			layer.MakeTiles(layer.GridTiles, opt)
		}
	}
}

// MakeTiles parses the tiles and loads the required images
func (layer *Layer) MakeTiles(tileLayer []*Tile, opt Options) {
	path := layer.Tileset
	if opt.Aseprite {
		path = strings.Replace(layer.Tileset, "aseprite", "png", 1)
	}
	img := newImageFromEmbed(opt.Root+"/"+path, opt.EmbedFolder)
	layer.Tiles = make(map[int]*ebiten.Image)
	for _, tileData := range tileLayer {
		if _, ok := layer.Tiles[tileData.Type]; !ok {
			tile := img.SubImage(
				image.Rect(
					tileData.Src[0],
					tileData.Src[1],
					tileData.Src[0]+layer.Size,
					tileData.Src[1]+layer.Size,
				),
			).(*ebiten.Image)
			layer.Tiles[tileData.Type] = tile
		}
	}
}

func (level *Level) MakeCollisions(opt Options) {
	for _, layer := range level.Layers {
		if len(layer.IntGrid) > 0 {
			layer.MakeCollisions(layer.IntGrid, opt)
		}
	}
}

func (layer *Layer) MakeCollisions(intGridTiles []*IntGridTile, opt Options) {
	for _, tile := range intGridTiles {
		if opt.CollidesWith != nil {
			if val, ok := opt.CollidesWith[tile.Value]; ok && val {
				var x, y float64
				y = math.Floor(float64(tile.Coord) / float64(layer.Width))
				x = float64(tile.Coord) - y*float64(layer.Width)
				opt.OnCollisionAdd(x*float64(layer.Size), y*float64(layer.Size), float64(layer.Size), float64(layer.Size))
			}
		}
	}
}

// newImageFromEmbed loads an image from an embedded file system
func newImageFromEmbed(path string, folder *embed.FS) *ebiten.Image {
	imgData, err := folder.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	imgFile, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(imgFile)
}

func (tile *Tile) Draw(screen *ebiten.Image, tiles map[int]*ebiten.Image) {
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(tile.Px[0]), float64(tile.Px[1]))
	screen.DrawImage(tiles[tile.Type], opt)
}

func (layer *Layer) Draw(screen *ebiten.Image) {
	for _, tile := range layer.AutoLayerTiles {
		tile.Draw(screen, layer.Tiles)
	}
	for _, tile := range layer.GridTiles {
		tile.Draw(screen, layer.Tiles)
	}
}

func (ldtkMap *Ldtk) Draw(screen *ebiten.Image) {
	for _, layer := range ldtkMap.Active.Layers {
		if len(layer.Tiles) > 0 {
			layer.Draw(screen)
		}
	}
}

func (ldtkMap *Ldtk) FindEntity(name string) *Entity {
	for _, layer := range ldtkMap.Active.Layers {
		if layer.Type == "Entities" {
			for _, entity := range layer.Entities {
				if entity.ID == name {
					return entity
				}
			}
		}
	}
	return nil
}
