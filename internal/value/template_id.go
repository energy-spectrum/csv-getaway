package value

import "strconv"

type TemplateID uint32

func (c *TemplateID) String() string {
	if c != nil {
		return strconv.Itoa(int(*c))
	}

	return ""
}
