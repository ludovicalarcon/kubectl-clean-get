package cmd

import (
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestReadFromFileAsInputIsNotDash(t *testing.T) {
	data := []byte(`foo: bar`)
	path := t.TempDir()
	Input = path + "/tmp.yaml"
	err := os.WriteFile(Input, data, 0644)

	if err != nil {
		t.Fatal("Can't write tempory file", err)
	}

	got, _ := readInput(Input)
	want := data

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestReadFromStdinAsInputIsNotDash(t *testing.T) {
	Input = "-"
	in := []byte("foo: bar")
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	_, err = w.Write(in)
	if err != nil {
		t.Error(err)
	}
	w.Close()

	// Restore stdin right after the test.
	defer func(v *os.File) { os.Stdin = v }(os.Stdin)
	os.Stdin = r

	got, _ := readInput(Input)
	want := []byte("foo: bar\n")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestReadFromStdinAnndHandleEOF(t *testing.T) {
	Input = "-"
	in := []byte("foo: bar\nEOF")
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	_, err = w.Write(in)
	if err != nil {
		t.Error(err)
	}
	w.Close()

	// Restore stdin right after the test.
	defer func(v *os.File) { os.Stdin = v }(os.Stdin)
	os.Stdin = r

	got, _ := readInput(Input)
	want := []byte("foo: bar\n")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestReadFromFileShouldHandleError(t *testing.T) {
	Input = "tmp.yaml"

	_, err := readInput(Input)

	if err == nil {
		t.Error("It should have asserted as input file doesn't exist.")
	}
}

func TestShouldCleanUidField(t *testing.T) {
	manifest := []byte(`foo: bar
uid: 123`)

	var data map[string]interface{}
	err := yaml.Unmarshal(manifest, &data)
	if err != nil {
		t.Fatal("Can't unmarshal yaml", err)
	}

	got, _ := cleanManifestInput(data)
	want := []byte("foo: bar\n")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestShouldHandleInputError(t *testing.T) {
	Input = "tmp.yaml"

	err := rootCmd.RunE(rootCmd, []string{})

	if err == nil {
		t.Error("It should have asserted as input file doesn't exist.")
	}
}

func TestShouldHandleInvalidYaml(t *testing.T) {
	Input = "tmp.yaml"
	data := []byte(`foo\nbar`)
	path := t.TempDir()
	Input = path + "/tmp.yaml"
	err := os.WriteFile(Input, data, 0644)
	if err != nil {
		t.Fatal("Can't write tempory file", err)
	}

	err = rootCmd.RunE(rootCmd, []string{})

	if err == nil {
		t.Error("It should have asserted as input yaml is invalid.")
	}
}

func TestShouldNotReturnErrorWhenEveyrthingOk(t *testing.T) {
	data := []byte(`foo: bar`)
	path := t.TempDir()
	Input = path + "/tmp.yaml"
	err := os.WriteFile(Input, data, 0644)

	if err != nil {
		t.Fatal("Can't write tempory file", err)
	}

	err = rootCmd.RunE(rootCmd, []string{})

	if err != nil {
		t.Errorf("It should not have asserted.")
	}
}