package utils

func ReverseStringArr(input []string) {
	for i, j := 0, len(input)-1; i < j; i, j = i+1, j-1 {
		input[i], input[j] = input[j], input[i]
	}
}

func SplitStringArrByLength(s string, n int) []string {
	if n <= 0 || len(s) == 0 {
			return nil
	}
	chunks := make([]string, 0, (len(s)+n-1)/n)
	for i := 0; i < len(s); i += n {
			end := i + n
			if end > len(s) {
					end = len(s)
			}
			chunks = append(chunks, s[i:end])
	}
	return chunks
}