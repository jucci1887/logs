/*
 Author: Kernel.Huang
 Mail: kernelman79@gmail.com
 File: config.go
 Date: 8/9/22 1:57 PM
*/
package logs

import (
	goToml "github.com/pelletier/go-toml"
	"log"
)

type TomlConfig struct {
	keyName    string
	value      interface{}
	Structured interface{}
	cfg        *goToml.Tree
}

var Toml = new(TomlConfig)

func (tf *TomlConfig) NewToml(dirname string, filename string) *TomlConfig {
	name := GetCustomConfigPath(dirname, filename)
	conf, err := goToml.LoadFile(name)

	if err != nil {
		log.Println("Load toml file error: ", err)
	}

	tf.cfg = conf
	return tf
}

// Example: result := Tome.NewToml(dirname, filename).Zone("zoneName").Get("key").To()
func (tf *TomlConfig) Zone(key string) *TomlConfig {
	tf.keyName = key
	return tf
}

// Example: result := Tome.NewToml(dirname, filename).Zone("zoneName").Get("key").To()
func (tf *TomlConfig) Get(key string) *TomlConfig {
	tf.keyName = tf.keyName + "." + key
	return tf
}

/**
 * Example: result := Tome.NewToml(dirname, filename).Zone("zoneName").Get("key").To()
 */
func (tf *TomlConfig) To() interface{} {
	return tf.cfg.Get(tf.keyName)
}

// Example: result := Tome.NewToml(dirname, filename).Zone("zoneName").Get("key").AtStr()
func (tf *TomlConfig) AtStr() string {
	tf.value = tf.cfg.Get(tf.keyName)
	return tf.value.(string)
}

// Example: result := Tome.NewToml(dirname, filename).Zone("zoneName").Get("key").AtInt()
func (tf *TomlConfig) AtInt() int {
	tf.value = tf.cfg.Get(tf.keyName)
	return tf.value.(int)
}

// Example: result := Tome.NewToml(dirname, filename).Zone("zoneName").Get("key").AtBool()
func (tf *TomlConfig) AtBool() bool {
	tf.value = tf.cfg.Get(tf.keyName)
	return tf.value.(bool)
}

// Example: result := Tome.NewToml(dirname, filename).Zone("zoneName").Fetch("key").ToStr()
func (tf *TomlConfig) Fetch(key string) *TomlConfig {
	tf.keyName = tf.keyName + "." + key
	tf.value = tf.cfg.Get(tf.keyName)
	return tf
}

// Example: result := Tome.NewToml(dirname, filename).Read("zoneName.key").ToStr() or ToInt()
func (tf *TomlConfig) Read(key string) *TomlConfig {
	tf.keyName = key
	tf.value = tf.To()
	return tf
}

// Example: result := Tome.NewToml(dirname, filename).Read("zoneName.key").ToStr()
func (tf *TomlConfig) ToStr() string {
	return tf.value.(string)
}

// Example: result := Tome.NewToml(dirname, filename).Read("zoneName.key").ToInt()
func (tf *TomlConfig) ToInt() int {
	return tf.value.(int)
}

// Example: result := Tome.NewToml(dirname, filename).Read("zoneName.key").ToBool()
func (tf *TomlConfig) ToBool() bool {
	return tf.value.(bool)
}
