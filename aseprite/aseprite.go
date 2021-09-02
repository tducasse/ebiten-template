package aseprite

import (
	"bytes"
	"embed"
	"encoding/json"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Frame struct {
	Image    *ebiten.Image
	Duration float64
}

type Tag struct {
	From int
	To   int
	Name string
}

type Animation struct {
	Tags            map[string]*Tag
	Frames          []*Frame
	CurrentTag      *Tag
	CurrentFrame    *Frame
	CurrentFrameIdx int
	FrameCounter    float64
	Options         *Options
	Image           *ebiten.Image
	OnLoopCallback  OnLoopType
	Instance        interface{}
}

type OnLoopType func(interface{}, *Animation)

type Options struct {
	// a pointer to the folder containing assets
	EmbedFolder *embed.FS
	// the root path to the aseprite assets (json and png)
	Root string
}

func Load(path string, opt *Options, tag string, instance interface{}) *Animation {
	data, err := opt.EmbedFolder.ReadFile(opt.Root + "/" + path)
	if err != nil {
		log.Fatal(err)
	}
	var anim Animation
	anim.Options = opt
	anim.Instance = instance
	var aseData map[string]interface{}
	err = json.Unmarshal(data, &aseData)
	if err != nil {
		log.Fatal(err)
	} else {
		anim.Init(aseData, tag)
	}
	return &anim
}

func (anim *Animation) Init(aseData map[string]interface{}, tag string) {
	frames := aseData["frames"].(map[string]interface{})
	meta := aseData["meta"].(map[string]interface{})
	anim.GetImage(meta)
	anim.GetTags(meta)
	anim.GetFrames(frames)
	anim.CurrentTag = anim.Tags[tag]
	anim.CurrentFrameIdx = anim.CurrentTag.From
}

func (anim *Animation) OnLoop(callback OnLoopType) {
	anim.OnLoopCallback = callback
}

func (anim *Animation) GetImage(meta map[string]interface{}) {
	imagePath := meta["image"].(string)
	anim.Image = newImageFromEmbed(anim.Options.Root+"/"+imagePath, anim.Options.EmbedFolder)
}

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

func (anim *Animation) GetTags(meta map[string]interface{}) {
	if anim.Tags == nil {
		anim.Tags = make(map[string]*Tag)
	}
	for _, val := range meta["frameTags"].([]interface{}) {
		tag := val.(map[string]interface{})
		anim.Tags[tag["name"].(string)] = &Tag{
			From: int(tag["from"].(float64)),
			To:   int(tag["to"].(float64)),
			Name: tag["name"].(string),
		}
	}
}

func (anim *Animation) GetFrames(frames map[string]interface{}) {
	for _, val := range frames {
		frameItem := val.(map[string]interface{})
		frame := frameItem["frame"].(map[string]interface{})
		x := int(frame["x"].(float64))
		y := int(frame["y"].(float64))
		w := int(frame["w"].(float64))
		h := int(frame["h"].(float64))
		duration := frameItem["duration"].(float64)
		anim.Frames = append(
			anim.Frames,
			&Frame{
				Image:    anim.Image.SubImage(image.Rect(x, y, x+w, y+h)).(*ebiten.Image),
				Duration: duration,
			})
	}
}

func (anim *Animation) Update() {
	if anim.CurrentFrame == nil {
		anim.CurrentFrame = anim.Frames[anim.CurrentFrameIdx]
		return
	}
	if anim.FrameCounter >= anim.CurrentFrame.Duration {
		newFrameIdx := anim.CurrentFrameIdx + 1
		if newFrameIdx > anim.CurrentTag.To {
			if anim.OnLoopCallback != nil {
				anim.OnLoopCallback(anim.Instance, anim)
			}
			newFrameIdx = anim.CurrentTag.From
		}
		anim.CurrentFrameIdx = newFrameIdx
		anim.FrameCounter = 0
	} else {
		anim.CurrentFrame = anim.Frames[anim.CurrentFrameIdx]
		anim.FrameCounter += (1.0 / 60.0) * 1000
	}
}

func (anim *Animation) Draw(x float64, y float64, screen *ebiten.Image) {
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(x, y)
	screen.DrawImage(anim.CurrentFrame.Image, opt)
}

func (anim *Animation) SetTag(tag string) {
	if anim.Tags[tag] == nil {
		log.Println("Could not find tag " + tag)
		return
	}
	if anim.CurrentTag.Name == tag {
		return
	}
	anim.CurrentTag = anim.Tags[tag]
	anim.CurrentFrameIdx = anim.CurrentTag.From
	anim.CurrentFrame = anim.Frames[anim.CurrentFrameIdx]
}
