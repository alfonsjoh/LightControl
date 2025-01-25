package Colors

type Color interface {
	GetColor() (string, error)
}

func ColorEquals(c1 Color, c2 Color) bool {
	c1Str, err1 := c1.GetColor()
	c2Str, err2 := c2.GetColor()
	return err1 == nil && err2 == nil && c1Str == c2Str
}
