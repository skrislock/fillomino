# fillomino

NEXT: Toothpaste Tube Force

Tests: TestSolveBoard2, estToothpasteFill_NegativeTest

Present Board 2

5 6 6 6 2 6 6 6 7 7
5 5 5 6 2 6 1 6 7 2
2 2 5 6 6 ! ! 6 7 2
5 5 3 _ _ 6 & 1 7 7
5 3 3 6 _ $ 7 7 _ 7
5 5 * 6 2 _ 6 7 7 1
4 3 3 6 2 5 _ 5 7 5
4 3 4 6 5 5 _ 5 5 5
4 2 4 6 5 2 5 2 2 1
4 2 4 4 5 2 5 5 5 5

6 is only one size now, and there is space up and to the right (&).  However, it could only grow one size.  It's only other option is to grow down, so it must go to $.  And when it goes to $, then the column of 6 in j=3 must then grow to the left *, as it cannot go up or to the right (logic exists, if otherOptions is properly excluding bad jumping writes).  Also, the 7 must then grow up to & (logic exists)

We are on a number, we have more than 1 possible write, or we have 1 possible write and otherOptions


countSoFar = count that number -- countContinuous
Create an "alreadyAnalyzed" map that can be reused

possibleWriteSpaceCount = Using the alreadyAnalyzed map, count the adjacent blanks for each possible write -- countContinuous
    * improvement: don't count anything next to something that is in the done set and is of equal value -- findAdjacentValues, remove anything from doneSet -- leave alreadyAnalyzed alone and reduce the count, since we need it for later
    * Question:  Where do we store these values
    * Question: Need to consider each possibleWrite as it's own possibility?  Or just take the smallest one(s) as forcers and the biggest count as the possible force?  Calculate each on it's own and the math will make it work?
        my 6 example: Go right.  My possible space = 1, count 1 + down 6, infinity > targetValue, don't write
                      Go down. My possible space = 6/infinity, count 1 plus right 4(1) < targetValue, so force down
if possibleWriteCount + countSoFar >= targetValue  no need to calculate otherOptions, as we are not in a position to force

otherOptionsSpaceCount -- using same alreadyAnalyzed counter? count the total in other options
    * improvement: same as above, minus equal value, done set, but need to get count right in combination with the possibleWriteSpaceCount

Idea:  OtherOptions and Possible Write are the same thing?


FUTURE BOARD 2 (board2KeyTest)
When the above logic works, the board will get to:
5 6 6 6 2 6 6 6 7 7
5 5 5 6 2 6 1 6 7 2
2 2 5 6 6 1 7 6 7 2
5 5 3 _ _ 6 7 1 7 7
5 3 3 6 _ 6 7 7 1 7
5 5 6 6 2 6 6 7 7 1
4 3 3 6 2 5 6 5 7 5
4 3 4 6 5 5 6 5 5 5
4 2 4 6 5 2 5 2 2 1
4 2 4 4 5 2 5 5 5 5

The system should write 1 (can't be a 2) in 4,4 but it thinks it may be able to write a 3, so it dosn't write anything.  We should eliminate 3 as a possible write becuase it is the value of the bound space and it is adjacent to the bound space

Present Board 5

4 4 4 4 _ _ _ _ 4 4
3 3 3 _ 4 4 _ _ 2 4
2 4 ! 3 3 2 2 _ 2 4
2 ! ! 3 4 4 _ 4 4 3
3 3 _ _ _ _ _ _ 1 3
3 4 4 _ 2 _ _ _ _ 3
1 2 2 3 2 4 _ _ 2 1
4 3 1 3 3 2 2 _ 2 3
4 3 3 _ 4 3 _ _ 3 3
4 4 _ _ _ _ _ 4 4 4

! 4 must be a square, as it can't leave beyond that square (jumping 4) and the size is four.  Diagonal must be OK at least

FUTURE FIVE TestBoard5KeySet

4 4 4 4 1 4 4 1 4 4
3 3 3 1 4 4 1 4 2 4
2 4 4 3 3 2 2 4 2 4
2 4 4 3 4 4 1 4 4 3
3 3 1 4 1 4 4 _ 1 3
3 4 4 4 2 _ _ _ _ 3
1 2 2 3 2 4 4 4 2 1
4 3 1 3 3 2 2 _ 2 3
4 3 3 _ 4 3 _ _ 3 3
4 4 _ _ _ _ _ 4 4 4

5,8 is the crux, it must be 5, can't be 1 2 3 (adjacent doneset), can't be 4 (can't grow anywhere), space is 5