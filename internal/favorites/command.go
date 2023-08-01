package favorites

import (
	"time"

	"github.com/gerladeno/favorites-mechanics/pkg/list"
)

type entry struct {
	id        int                       `yaml:"id"`
	name      string                    `yaml:"name"`
	exec      string                    `yaml:"exec"`
	entries   *list.DeLinkedList[entry] `yaml:"entries"`
	isDir     bool                      `yaml:"isDir"`
	createdAt time.Time                 `yaml:"createdAt"`
	updatedAt time.Time                 `yaml:"updatedAt"`
}
