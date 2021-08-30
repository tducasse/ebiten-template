package ldtk

import (
	"bytes"
	"embed"
	"encoding/json"
	"image"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
	Px   [2]int `json:"px"`
	Src  [2]int `json:"src"`
	Type int    `json:"t"`
	Flip int    `json:"f"`
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
	ID             string    `json:"__identifier"`
	Type           string    `json:"__type"`
	Entities       []*Entity `json:"entityInstances"`
	Size           int       `json:"__gridSize"`
	Width          int       `json:"cWid"`
	Height         int       `json:"cHei"`
	AutoLayerTiles []*Tile   `json:"autoLayerTiles"`
	GridTiles      []*Tile   `json:"gridTiles"`
	Tileset        string    `json:"__tilesetRelPath"`
	IntGridCSV     []int     `json:"intGridCSV"`
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
	EmbedFolder *embed.FS
	Aseprite    bool
	FilePrefix  string
}

func Load(data []byte, opt *Options) (*Ldtk, error) {
	var ldtkMap Ldtk
	ldtkMap.Options = opt
	err := json.Unmarshal(data, &ldtkMap)
	if err != nil {
		return nil, err
	}
	for _, level := range ldtkMap.Levels {
		level.PrepareLevel(opt)
	}
	ldtkMap.Active = ldtkMap.Levels[0]
	return &ldtkMap, nil
}

// PrepareLevel populates the layer.Tiles prop if required
func (level *Level) PrepareLevel(opt *Options) {
	for _, layer := range level.Layers {
		if len(layer.AutoLayerTiles) > 0 {
			layer.MakeTiles(&layer.AutoLayerTiles, opt)
		} else if len(layer.GridTiles) > 0 {
			layer.MakeTiles(&layer.GridTiles, opt)
		}
	}
}

// MakeTiles parses the tiles and loads the required images
func (layer *Layer) MakeTiles(tileLayer *[]*Tile, opt *Options) {
	path := layer.Tileset
	if opt.Aseprite {
		path = strings.Replace(layer.Tileset, "aseprite", "png", 1)
	}
	img := newImageFromEmbed(opt.FilePrefix+"/"+path, opt.EmbedFolder)
	layer.Tiles = make(map[int]*ebiten.Image)
	for _, tileData := range *tileLayer {
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

func (tile *Tile) Draw(screen *ebiten.Image, tiles *map[int]*ebiten.Image) {
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(tile.Px[0]), float64(tile.Px[1]))
	screen.DrawImage((*tiles)[tile.Type], opt)
}

func (layer *Layer) Draw(screen *ebiten.Image) {
	for _, tile := range layer.AutoLayerTiles {
		tile.Draw(screen, &layer.Tiles)
	}
	for _, tile := range layer.GridTiles {
		tile.Draw(screen, &layer.Tiles)
	}
}

func (ldtkMap *Ldtk) Draw(screen *ebiten.Image) {
	for _, layer := range ldtkMap.Active.Layers {
		if len(layer.Tiles) > 0 {
			layer.Draw(screen)
		}
	}
}

// for _, field := range levels.Levels[0].Layers[0].Entities[0].Fields {
// 	if strings.HasPrefix(field.Type, "Array") {
// 		myField, ok := field.Value.([]interface{})
// 		if !ok {
// 			log.Printf("nok")
// 		} else {
// 			log.Println(myField)
// 		}
// 	} else {
// 		log.Println(field.Value)
// 	}
// }
