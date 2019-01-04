package heldiamgo

import (
	"image"
	"image/color"
)

//绘图

//圆形 https://blog.golang.org/go-imagedraw-package
type Circle struct {
	p image.Point
	r int
}

//新建一个圆形，p为圆形中心点. r为半径
func NewCircle(p image.Point, r int) *Circle {
	return &Circle{
		p: p,
		r: r,
	}
}

func (c *Circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *Circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *Circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{A: 255}
	}
	return color.Alpha{A: 0}
}

//椭圆形
type Ellipse struct {
	rect   image.Rectangle
	radius int
	width  int
	height int
}

func (c *Ellipse) SetBounds(rect image.Rectangle) {
	c.rect = rect
	c.width = rect.Size().X
	c.height = rect.Size().Y
}

func (c *Ellipse) SetRadius(radius int) {
	c.radius = radius
}

func (c *Ellipse) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *Ellipse) Bounds() image.Rectangle {
	return c.rect
}

func (c *Ellipse) At(x, y int) color.Color {
	if x < c.radius && y < c.radius {
		xx, yy, rr := float64(x-c.radius)+0.5, float64(y-c.radius)+0.5, float64(c.radius)
		if xx*xx+yy*yy < rr*rr {
			return color.Alpha{A: 255}
		}
		return color.Alpha{A: 0}
	} else if x > c.width-c.radius && y < c.radius {
		xx, yy, rr := float64(x+c.radius-c.width)+0.5, float64(y-c.radius)+0.5, float64(c.radius)
		if xx*xx+yy*yy < rr*rr {
			return color.Alpha{A: 255}
		}
		return color.Alpha{A: 0}
	} else if x > c.width-c.radius && y > c.height-c.radius {
		xx, yy, rr := float64(x+c.radius-c.width)+0.5, float64(y+c.radius-c.height)+0.5, float64(c.radius)
		if xx*xx+yy*yy < rr*rr {
			return color.Alpha{A: 255}
		}
		return color.Alpha{A: 0}
	} else if x < c.radius && y > c.height-c.radius {
		xx, yy, rr := float64(x-c.radius)+0.5, float64(y+c.radius-c.height)+0.5, float64(c.radius)
		if xx*xx+yy*yy < rr*rr {
			return color.Alpha{A: 255}
		}
		return color.Alpha{A: 0}
	}
	return color.Alpha{A: 255}
}
