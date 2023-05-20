package renders

import "fmt"

func formatAmount(a int) string {
	f := float32(a) / float32(100)
	return fmt.Sprintf("$%.2f", f)
}
