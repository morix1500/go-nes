package ui

import (
	"errors"
	"image"
	"image/draw"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var errUnsupportedStride = errors.New("unsupported stride, only 32-bit colors supported")

var errTextureNotBound = errors.New("texture not bound")

const (
	textureSize = 4096
	textureDim  = textureSize / 256
)

type Texture struct {
	handle  uint32
	target  uint32 // same target as gl.BindTexture(<this param>, ...)
	texUnit uint32 // Texture unit that is currently bound to ex: gl.TEXTURE0
}

func NewTexture(img image.Image, wrapR, wrapS int32) (*Texture, error) {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	if rgba.Stride != rgba.Rect.Size().X*4 { // TODO-cs: why?
		return nil, errUnsupportedStride
	}

	var handle uint32
	gl.GenTextures(1, &handle)

	target := uint32(gl.TEXTURE_2D)
	//internalFmt := int32(gl.SRGB_ALPHA)
	format := uint32(gl.RGBA)
	width := int32(rgba.Rect.Size().X)
	height := int32(rgba.Rect.Size().Y)
	pixType := uint32(gl.UNSIGNED_BYTE)
	dataPtr := gl.Ptr(rgba.Pix)

	texture := Texture{
		handle: handle,
		target: target,
	}

	texture.Bind(gl.TEXTURE0)
	defer texture.UnBind()

	// set the texture wrapping/filtering options (applies to current bound texture obj)
	// TODO-cs
	gl.TexParameteri(texture.target, gl.TEXTURE_WRAP_R, wrapR)
	gl.TexParameteri(texture.target, gl.TEXTURE_WRAP_S, wrapS)
	gl.TexParameteri(texture.target, gl.TEXTURE_MIN_FILTER, gl.LINEAR) // minification filter
	gl.TexParameteri(texture.target, gl.TEXTURE_MAG_FILTER, gl.LINEAR) // magnification filter

	//gl.TexImage2D(target, 0, gl.RGBA, width, height, 0, format, pixType, nil)
	gl.TexImage2D(target, 0, gl.RGBA, width, height, 0, format, pixType, dataPtr)

	gl.GenerateMipmap(texture.handle)

	return &texture, nil
}

func (tex *Texture) Bind(texUnit uint32) {
	gl.ActiveTexture(texUnit)
	gl.BindTexture(tex.target, tex.handle)
	tex.texUnit = texUnit
}

func (tex *Texture) UnBind() {
	tex.texUnit = 0
	gl.BindTexture(tex.target, 0)
}

func (tex *Texture) SetUniform(uniformLoc int32) error {
	if tex.texUnit == 0 {
		return errTextureNotBound
	}
	gl.Uniform1i(uniformLoc, int32(tex.texUnit-gl.TEXTURE0))
	return nil
}
