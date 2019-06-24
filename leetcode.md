1. 只出现一次的数字
   异或解法,erlang版本:
   single_number(NumberList) ->
     lists:foldl(fun(N,R) -> N bxor R end, 0, NumberList). 
   借助额外数组,golang版本:
   func singleNumber(nums []int) int {
	  numberMap := make(map[int]int)

 	  for i := 0; i < len(nums); i++ {
 		 count, ok := numberMap[nums[i]]
 		 if ok {
 			numberMap[nums[i]] = count + 1
 		 } else {
 			numberMap[nums[i]] = 1
 		 }
 	  }

     var result int
 	  for i := range(numberMap) {
 		 if numberMap[i] == 1 {
 			 result = i
 		 }
 	 }
    
     return result
   }
2. 求众数
   常规解法,golang版本:
   func majorityElement(nums []int) int {
    numberMap := make(map[int]int)

    for i := 0; i < len(nums); i++ {
 		count, ok := numberMap[nums[i]]
 		if ok {
 			numberMap[nums[i]] = count + 1
 		} else {
 			numberMap[nums[i]] = 1
 		}
 	}

    result, count := 0 , 0
    for i := range(numberMap) {
    	if numberMap[i] >= len(nums)/2 && numberMap[i] >= count {
    		result, count = i, numberMap[i]
    	}
     }

     return result
  }
  摩尔投票法(Moore Voting),前提是必定存在出现次数过半的数,
  摩尔投票法先假定第一个元素过半数,计数器设为1,比较下一个数和次数是否相同,若相同则计数器+1,反之-1,之后判断计数器是否为0,若为0,则将下一个数设为候选
  众数,以此类推,直至遍历完整个数组,当前候选即为要求的众数
  摩尔投票参考:https://www.zhihu.com/question/49973163
3. 搜索二维矩阵ii,由于行列均满足有序,最开始的想法是每一行执行二分查找,或者将二维数据归并排序,合并成一维数组,再执行二分查找,时间复杂度不符合'高效'的
   要求,基于从左到右,从上到下升序的前提,可以从左下角开始比较,若当前值等于目标值,则返回true,若当前值小于目标值,则右移,反之则左移,直到矩阵的右上角
   以下为golang版本的实现:
   func searchMatrix(matrix [][]int, target int) bool {
	for row, col := len(matrix) - 1, 0; row >= 0 && col <= len(matrix[0]) -1; {
		if matrix[row][col] == target {
			return true
		} else if matrix[row][col] > target {
			row -= 1
		} else {
			col += 1
		}
	}

	return false
    }
4. 两数之和,常规思路是双层循环,找到目标之后直接返回,时间复杂度高,其次是借助使用hash结构的数据类型缓存中间计算结果,一下为使用map的golang版本:
	func twoSum(nums []int, target int) []int {
    var result []int
    var diffs map[int]int
    diffs = make(map[int]int)

    for i := 0; i < len(nums); i++ {
    	diff := target - nums[i]
    	if index, ok := diffs[diff]; ok {
    		if index > i {
    			result = append(result,i)
    			result = append(result,index)
    			return result
    		} else {
    			result = append(result,index)
    			result = append(result,i)
    			return result
    		}
    	}

    	diffs[nums[i]] = i
     }
    
      return result
   }
