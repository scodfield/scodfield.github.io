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
5. 两数相加,常规思路是先判断两个链表是否为nil,循环两个不为空的链表,最后在判断进位是否大于0,大于0则需要多一个节点,另一种思路是用递归来实现
   以下代码为golang版的递归实现:
   func addTwoNumbers2(l1 *ListNode, l2 *ListNode) *ListNode { 
    return addtwo(l1,l2,0)
   }

   func addtwo(l1 *ListNode, l2 *ListNode, carry int) *ListNode {
	x, y := 0, 0
	if l1 != nil {
		x = l1.Val
	}

	if l2 != nil {
		y = l2.Val
	}

	sum := x + y + carry

       var newNode ListNode
       newNode.Val = sum % 10
       carry = sum / 10
       if l1.Next == nil || l2.Next == nil {
    	   if l1.Next != nil {
    		var l2Node ListNode
    		l2Node.Val = 0
    		l2.Next = &l2Node
    		newNode.Next = addtwo(l1.Next,l2.Next,carry)
    	   } else if l2.Next != nil {
    		var l1Node ListNode
    		l1Node.Val = 0
    		l1.Next = &l1Node
    		newNode.Next = addtwo(l1.Next,l2.Next,carry)
    	   } else if carry > 0 {
    		var lastNode ListNode
    		lastNode.Val = carry
    		newNode.Next = &lastNode
    	   }
       } else {
    	  newNode.Next = addtwo(l1.Next,l2.Next,carry)
       }

      return &newNode
   }
6. 无重复字符的最长子串,常规思路还是双层循环,统计每个字符可以构造的最长无重复子串,最后返回最长子串的长度即可,贴一个自己的实现:
    func lengthOfLongestSubstring(s string) int {
    var result int = 0

    if s == "" {
    	return 0
    } else if len(s) == 1 {
    	return 1
    }

    var count int = 1
    for i := 0; i < len(s) - 1; i++ {
    	temp := make(map[uint8]int)
    	temp[s[i]] = 1
    	count = 1
    	for j := i+1; j < len(s); j++ {
    		_, ok := temp[s[j]]
    		if ok {
    			break
    		} else {
    			count ++
    			temp[s[j]] = 1
    		}
    	 }

    	 if result < count {
    		result = count
    	 }
      }

      return result
   }
     上述代码提交后,发现不关是时间还是空间都很垃圾,参考了评论,发现有更好的思路,代码如下:
     func lengthOfLongestSubstring(s string) int {
	i, j, max := 0, 0, 0

	for ; i < len(s); i++ {
		for x := j; x < i; x++ {
			if s[x] == s[i] {
				j = x + 1
				break
			}
		}
		fmt.Println("i: ",i,", j: ",j)
		if i-j+1 > max {
			max = i - j + 1
		}
	}
	return max
     }
     由代码可知解题思路为求当前下标到最大下标之间的不重复子串,一旦出现重复字符,则当前子串查找结束,判断是否是当前最长,更新max的值,同时更改当前下标
     为当前重复字符的下一个下标值,i,j,max初始值及相关计算方式可以兼顾到s的不同情况(为空或者只有一个字符等)
7. 寻找两个有序数组的中位数,题目要去log(m+n),目前只实现了m+n的时间复杂度,首先合并两个有序数组,再根据奇偶计算中位数,实现如下:
   func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
    // nums3 := merge(nums1,len(nums1),nums2,len(nums2))
    nums3 := merge_sorted(nums1,nums2)
    fmt.Println("median nums3",nums3)
    if len(nums3) % 2 == 0 {
    		return float64(nums3[len(nums3)/2] + nums3[len(nums3)/2-1]) / float64(2)
    	} else {
    		return float64(nums3[len(nums3)/2])
    	}
   }

   func merge_sorted(nums1 []int, nums2 []int) []int {
	var nums3 []int
	i, j := 0, 0

	for i < len(nums1) && j < len(nums2) {
		if nums1[i] < nums2[j] {
			nums3 = append(nums3,nums1[i])
			i++
		} else if nums1[i] > nums2[j] {
			nums3 = append(nums3,nums2[j])
			j++
		} else {
			nums3 = append(nums3,nums2[j])
			nums3 = append(nums3,nums2[j])
			i++
			j++
		}
	}

	for ; i < len(nums1); i++ {
		nums3 = append(nums3,nums1[i])
	}
	for ; j < len(nums2); j++ {
		nums3 = append(nums3,nums2[j])
	}

	return nums3
    }
8. 二叉树的中序遍历,常规实现为递归,但时间和空间复杂度均为O(n),中序遍历的一个非常经典的算法是莫里斯(Morris)遍历,它可以做到时间O(n),空间O(1)
   莫里斯算法的关键在于当前节点的左子树中的最右节点,将该左右节点的右空指针指向当前节点,用于遍历时候的回退,在中序遍历序列中,该最右节点刚好是当前
   节点的前驱节点
   Morris参考:https://blog.csdn.net/koalacoco/article/details/60464975
