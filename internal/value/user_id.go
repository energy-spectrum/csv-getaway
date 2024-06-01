package value

import "strconv"

type UserID uint32

func (c *UserID) String() string {
	if c != nil {
		return strconv.Itoa(int(*c))
	}

	return ""
}
