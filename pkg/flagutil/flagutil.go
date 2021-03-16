package flagutil

import (
	"flag"
	"fmt"
	"strconv"

	"gopkg.in/yaml.v3"
)

func ParseYAML(fs *flag.FlagSet, data []byte) error {
	var m map[string]interface{}
	err := yaml.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	fmt.Println("here", m)

	return setup(fs, "", m)
}

func setup(fs *flag.FlagSet, key string, value interface{}) error {
	if value == nil {
		return nil
	}

	switch x := value.(type) {
	case map[string]interface{}:
		for k, v := range x {
			err := setup(fs, Join(key, k, "."), v)
			if err != nil {
				return err
			}
		}

	default:
		str, err := stringify(x)
		if err != nil {
			return fmt.Errorf("option %q: %w", key, err)
		}

		err = fs.Set(key, str)
		if err != nil {
			return fmt.Errorf("option %q: %w", key, err)
		}
	}

	return nil
}

func Join(a, b, sep string) string {
	if a == "" {
		return b
	}
	return a + sep + b
}

func stringify(x interface{}) (string, error) {
	switch v := x.(type) {
	case string:
		return v, nil

	case
		bool,
		int,
		int8,
		int16,
		int32,
		int64,
		uint,
		uint8,
		uint16,
		uint32,
		uint64,
		uintptr:

		return fmt.Sprint(v), nil

	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil

	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil

	default:
		return "", fmt.Errorf("cannot convert %T to string", v)
	}
}
