package hardware

import (
	"fmt"
	"io"
)

// MachineReader reads single Machine configuration at a time. When there are no more Machine entries
// to be read, Read() returns io.EOF.
type MachineReader interface {
	Read() (Machine, error)
}

// MachineWriter writes Machine entries.
type MachineWriter interface {
	Write(Machine) error
}

// MachineValidator validates an instance of Machine.
type MachineValidator interface {
	Validate(Machine) error
}

// TranslateAll reads entries 1 at a time from reader and writes them to writer. When reader returns io.EOF,
// TranslateAll returns nil. Failure to return io.EOF from reader will result in an infinite loop.
func TranslateAll(reader MachineReader, writer MachineWriter, validator MachineValidator) error {
	for {
		fmt.Println("TranslateAll: reading next machine...")
		err := Translate(reader, writer, validator)

		if err == io.EOF {
			fmt.Println("TranslateAll: reached end of input.")
			return nil
		}

		if err != nil {
			fmt.Printf("TranslateAll: error occurred: %v\n", err)
			return err
		}

		fmt.Println("TranslateAll: successfully translated machine.")
	}
}

// Translate reads 1 entry from reader and writes it to writer. When reader returns io.EOF Translate
// returns io.EOF to the caller.
func Translate(reader MachineReader, writer MachineWriter, validator MachineValidator) error {
	fmt.Println("Translate: reading machine...")
	machine, err := reader.Read()
	if err == io.EOF {
		fmt.Println("Translate: reached EOF")
		return err
	}

	if err != nil {
		fmt.Printf("Translate: read error: %v\n", err)
		return fmt.Errorf("read: invalid hardware: %v", err)
	}

	fmt.Printf("Translate: machine read: %+v\n", machine)

	fmt.Println("Translate: validating machine...")
	if err := validator.Validate(machine); err != nil {
		fmt.Printf("Translate: validation failed: %v\n", err)
		return err
	}
	fmt.Println("Translate: validation succeeded.")

	fmt.Println("Translate: writing machine...")
	if err := writer.Write(machine); err != nil {
		fmt.Printf("Translate: write error: %v\n", err)
		return fmt.Errorf("write: %v", err)
	}
	fmt.Println("Translate: write succeeded.")

	return nil
}
