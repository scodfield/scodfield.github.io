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
