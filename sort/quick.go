package sort

func QuickSort(nums []int) []int {
	left, right := 0, len(nums)-1
	if len(nums) == 0 || left == right {
		return nums
	}
	p := nums[left]
	for left < right {
		for right > left && nums[right] > p {
			right--
		}
		if right > left {
			nums[left], nums[right] = nums[right], nums[left]
			left++
		}
		for left < right && nums[left] < p {
			left++
		}
		if right > left {
			nums[right], nums[left] = nums[left], nums[right]
			right--
		}
	}
	l := QuickSort(nums[:left])
	r := QuickSort(nums[left+1:])
	return append(l, append([]int{p}, r...)...)
}
