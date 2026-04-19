package core

import (
	"strings"

	"github.com/krozhkov/mgml/core/helpers"
)

type MJMLBuilder struct {
	strings.Builder
}

func (b *MJMLBuilder) WriteConditionalTag(s string, negation bool) (int, error) {
	var len int
	if i, err := helpers.ConditionalTagStart(b, negation); err != nil {
		return 0, err
	} else {
		len += i
	}

	if i, err := b.WriteString(s); err != nil {
		return 0, err
	} else {
		len += i
	}

	if i, err := helpers.ConditionalTagEnd(b, negation); err != nil {
		return 0, err
	} else {
		len += i
	}

	return len, nil
}

func (b *MJMLBuilder) WriteMsoConditionalTag(s string, negation bool) (int, error) {
	var len int
	if i, err := helpers.MsoConditionalTagStart(b, negation); err != nil {
		return 0, err
	} else {
		len += i
	}

	if i, err := b.WriteString(s); err != nil {
		return 0, err
	} else {
		len += i
	}

	if i, err := helpers.ConditionalTagEnd(b, negation); err != nil {
		return 0, err
	} else {
		len += i
	}

	return len, nil
}
