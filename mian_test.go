package main

import (
    "os"
    "testing"
)

func TestMain(m *testing.M) {
    os.Args = []string{"command", "-f", "fdse.cvs"}
    // code := m.Run()
    main()
    // os.Exit(code)
}
