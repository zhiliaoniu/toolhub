// qsort.go
package qsort

func quickSort2(values []int, left, right int) {
	if values == nil || left >= right {
		return
	}
	mid_data := values[left]
	begin, end := left, right
	for begin < end {
		for end > begin && values[end] >= mid_data {
			end--
		}
		values[begin] = values[end]
		for begin < end && values[begin] <= mid_data {
			begin++
		}
		values[end] = values[begin]
	}
	values[begin] = mid_data
	quickSort2(values, left, begin-1)
	quickSort2(values, begin+1, right)
}

func quickSort(values []int, left, right int) {
	temp := values[left]
	p := left
	i, j := left, right
	for i <= j {
		for j >= p && values[j] >= temp {
			j--
		}
		if j >= p {
			values[p] = values[j]
			p = j
		}
		if values[i] <= temp && i <= p {
			i++
		}
		if i <= p {
			values[p] = values[i]
			p = i
		}
	}
	values[p] = temp
	if p-left > 1 {
		quickSort(values, left, p-1)
	}
	if right-p > 1 {
		quickSort(values, p+1, right)
	}
}
func QuickSort(values []int) {
	quickSort2(values, 0, len(values)-1)
}
