package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Coordinate struct {
	I int
	J int
}

func solve(board [][]string) [][]string {
	h := len(board)
	w := len(board[0])

	// TODO Maybe doneSet can just be an array rather than a map, as the 'true' seems to be irrelevant?
	// make a set to keep track of what sections are DONE and can be eliminated as possibly growing
	doneSet := make(map[Coordinate]bool)
	doneSetPreviousSize := -1

	// TODO: remove the safety (20)
	for x := 0; x < 20; x++ {
		// TODO put this in a tight little function with doneSet, doneSetPreviousSize as pointers
		if len(doneSet) == h*w {
			fmt.Println("Solved!")
			break
		} else if len(doneSet) == doneSetPreviousSize {
			fmt.Println("Stuck... maybe?") // if there are calculations that require an additional iteration to have an effect, let it go once more?
			// TODO: not really good to stop here, just because a section is not done, it may be growing.
			// break
		} else {
			doneSetPreviousSize = len(doneSet)
			fmt.Println("DoneSet: ", doneSet)
		}

		// loop through and analyze each cell
		for i := 0; i < h; i++ {
			for j := 0; j < w; j++ {
				if doneSet[Coordinate{i, j}] {
					continue
				}
				if board[i][j] == "_" {
					// SUB SOLUTION:  SEED

					// check surrounding cells, and if they are all done, we can start to fill
					// in blanks.  Blank blocks larger than size 1 may take extra iteration to
					// determine fill value, eg. a blank size of 3 might be 1 block of 3 or 1 block of
					// one plus one block of two
					sizeOfBoundedBlank := analyzeBlankSpace(i, j, -1, nil, doneSet, &board)
					if sizeOfBoundedBlank == -1 {
						continue
					} else if sizeOfBoundedBlank == 1 || sizeOfBoundedBlank == 2 {
						// if it is size 2, then this value must be 2, as we can't have
						// two 1's together.  So we can seed this with a 2 the other space
						// will be solved by the rest of the logic
						board[i][j] = strconv.Itoa(sizeOfBoundedBlank)
					} else { // we need to see if we can seed this space with a value
						// for now...
						seedValue := determinePossibleSeedValue(i, j, sizeOfBoundedBlank, doneSet, &board)
						if seedValue != -1 {
							board[i][j] = strconv.Itoa(seedValue)
						}
					}
					continue
				} else {
					targetValue, err := strconv.Atoi(board[i][j])
					if err != nil {
						fmt.Println("error")
					}

					thisCount := countContinuous(i, j, 0, nil, &board)
					if thisCount == targetValue {
						// this could be made to be more efficient, but this should work, the iteration after
						// the section is complete
						// CONVERSELY, maybe leave as is to allow seed values to grow, and naturally be set to doneSet
						// when they are finally done
						doneSet[Coordinate{i, j}] = true
						continue
					} else if thisCount > targetValue {
						fmt.Printf("error at %v, %v, value is %v, count is %v\n", i, j, targetValue, thisCount)
					} else {
						alreadyBoundAnalyzed := make(map[Coordinate]bool)
						// todo maybe make this into a function that returns a map

						// identify possible writing directions
						possibleWrites := findAdjacentValues(i, j, &board, "_")
						// NEXT if joining to the one of possible writes would make it over the targetValue, eliminate it from the set

						// NEXT if joining to the one of possible writes would make it over the targetValue, eliminate it from the set
						// remove possible moves that will result in a too long section
						// TODO Maybe need to move or reshape this in next sub-solutions
						badJumpingWrites := identifyBadJumpingWrites(i, j, thisCount, &possibleWrites, &board, board[i][j])

						for k := range badJumpingWrites {
							alreadyBoundAnalyzed[Coordinate{k.I, k.J}] = true
							delete(possibleWrites, Coordinate{k.I, k.J})
						}

						if len(possibleWrites) == 0 {
							continue
						}

						// need to make sure that this is the only valid way this block can grow
						// so we check if other ways to grow this block exist
						hasOtherOption, otherOptionBoundValue := otherOptionExists(i, j, thisCount, targetValue, nil, alreadyBoundAnalyzed, doneSet, &board)

						// if just one, write
						if len(possibleWrites) == 1 && !hasOtherOption {
							for k := range possibleWrites { // TODO if this is an array or slice maybe this can be cleaner
								board[k.I][k.J] = board[i][j]
							}
						} else {
							// SUB SOLUTION: BOUND FORCE
							// This cell is completed surrounded by done cells or the wall, and the empty
							// space held within is equal to the cells value.  So all the empty spaces
							// must be this value, so write one of them here.  We know the surrounding cells
							// don't want to come in here because they are all in the doneSet
							alreadyBoundAnalyzed[Coordinate{i, j}] = true
							myBoundCount := 0
							for k := range possibleWrites {
								thisBoundCount := analyzeBlankSpace(k.I, k.J, targetValue, alreadyBoundAnalyzed, doneSet, &board)
								if thisBoundCount == -1 {
									myBoundCount = -1
									break
								} else {
									myBoundCount += thisBoundCount
								}
							}

							if (hasOtherOption && otherOptionBoundValue > 0) || len(possibleWrites) > 1 {
								if hasOtherOption && otherOptionBoundValue > 0 { // TODO is this supposed to be -1?
									myBoundCount += otherOptionBoundValue
								}
								if myBoundCount+thisCount == targetValue {
									for k := range possibleWrites { // TODO if this is an array or slice maybe this can be cleaner
										board[k.I][k.J] = board[i][j]
									}
									continue
								}
							}

							// we are not bound, but still might be in a positition to force
							if len(possibleWrites) > 1 || hasOtherOption { // hasOtherOptions infers len(possibleWrites) == 1
								alreadyForceAnalyzed := make(map[Coordinate]bool)

								for k := range possibleWrites { // TODO if this is an array or slice maybe this can be cleaner
									// consider this write (board[k.I][k.J]) for a force.
									mySpaceCount := countContinuous(k.I, k.J, 1, alreadyForceAnalyzed, &board)
									if mySpaceCount+thisCount <= targetValue {
										fmt.Println("Working On Force")
									}
								}
							}

						}
					}
				}
			}
		}

		// write out the board
		fmt.Println("after ", x+1, "iterations")
		printBoard(board)
	}

	return board
}

// todo make alreadyCounted a pointer
func countContinuous(i, j, counter int, alreadyCounted map[Coordinate]bool, board *[][]string) int {
	if alreadyCounted == nil {
		alreadyCounted = make(map[Coordinate]bool)
	}
	if alreadyCounted[Coordinate{i, j}] {
		return counter
	} else {
		// add this coordinate to the list
		counter++
		alreadyCounted[Coordinate{i, j}] = true
	}

	adjacentSameValues := findAdjacentValues(i, j, board, (*board)[i][j])

	for k := range adjacentSameValues {
		counter = countContinuous(k.I, k.J, counter, alreadyCounted, board)
	}

	return counter
}

func identifyBadJumpingWrites(i, j, countSoFar int, possibleWrites *map[Coordinate]bool, board *[][]string, value string) map[Coordinate]bool {
	targetValue, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println("error")
	}

	// I don't think it's safe to alter a set as I am ranging through it so I'll make a set to delete
	badJumpingWrites := make(map[Coordinate]bool)

	for y := range *possibleWrites {
		// find same value
		skipToSameValueCells := findAdjacentValues(y.I, y.J, board, value)
		for z := range skipToSameValueCells {
			if !(z.I == i && z.J == j) && countSoFar+countContinuous(z.I, z.J, 0, nil, board)+1 > targetValue {
				badJumpingWrites[Coordinate{y.I, y.J}] = true
			}
		}
	}

	return badJumpingWrites
}

// check if the solution block that you are presently in has other options that it might grow
func otherOptionExists(i, j, thisCount, targetValue int, alreadyCounted, alreadyBoundAnalyzed, doneSet map[Coordinate]bool, board *[][]string) (bool, int) {
	if alreadyCounted == nil {
		alreadyCounted = make(map[Coordinate]bool)
		// we dont want to count the one move we are thinking to make
		alreadyCounted[Coordinate{i, j}] = true
	}

	if alreadyBoundAnalyzed == nil {
		alreadyBoundAnalyzed = make(map[Coordinate]bool)
	}

	myBoundCount := 0

	adjacentBlankValues := findAdjacentValues(i, j, board, "_")
	if len(adjacentBlankValues) > 0 {
		// NEXT if joining to the one of possible writes would make it over the targetValue, eliminate it from the set
		// remove possible moves that will result in a too long section
		// TODO: Do we need to add these bad jumping writes to the alreadyBoundAnalyzed Coordinate map? -- testing now
		badJumpingWrites := identifyBadJumpingWrites(i, j, thisCount, &adjacentBlankValues, board, (*board)[i][j])
		for k := range badJumpingWrites {
			alreadyBoundAnalyzed[Coordinate{k.I, k.J}] = true
			delete(adjacentBlankValues, Coordinate{k.I, k.J})
		}
	}
	if !alreadyCounted[Coordinate{i, j}] && len(adjacentBlankValues) > 0 {
		alreadyBoundAnalyzed[Coordinate{i, j}] = true
		for abv := range adjacentBlankValues {
			// go through each adjacent blank value, see if it is bound
			thisBoundCount := analyzeBlankSpace(abv.I, abv.J, targetValue, alreadyBoundAnalyzed, doneSet, board)
			if thisBoundCount == -1 {
				// if it is not bound, return true, -1
				return true, -1
			} else {
				// if everything is bound, count up all the bounds, and carry on to the next adjacent targetValue spaces
				myBoundCount += thisBoundCount
			}
		}
	}

	alreadyCounted[Coordinate{i, j}] = true

	// check my adjacent cells
	adjacentSameValues := findAdjacentValues(i, j, board, (*board)[i][j])

	for k := range adjacentSameValues {
		if !alreadyCounted[Coordinate{k.I, k.J}] {
			hasOtherOption, otherOptionBoundValue := otherOptionExists(k.I, k.J, thisCount, targetValue,
				alreadyCounted, alreadyBoundAnalyzed, doneSet, board)
			if hasOtherOption {
				if otherOptionBoundValue == -1 {
					return true, -1
				} else {
					myBoundCount += otherOptionBoundValue
				}
			}
		}
	}

	if myBoundCount > 0 {
		return true, myBoundCount
	} else {
		return false, -1
	}
}

// returns a set of values adjacent to a given i, j that matches the value parameter
// TODO Simlar to doneSet, can the return value of this be changed to an array rather than
// a Coordinate - boolean map?  The boolean part seems to be irrelevant?
func findAdjacentValues(i, j int, board *[][]string, value string) map[Coordinate]bool {
	h := len(*board)
	w := len((*board)[0])
	adjacentValueCoordinates := make(map[Coordinate]bool)

	// check up
	if i > 0 && value == (*board)[i-1][j] {
		adjacentValueCoordinates[Coordinate{i - 1, j}] = true
	}

	// check down
	if i < (h-1) && value == (*board)[i+1][j] {
		adjacentValueCoordinates[Coordinate{i + 1, j}] = true
	}

	// check left
	if j > 0 && value == (*board)[i][j-1] {
		adjacentValueCoordinates[Coordinate{i, j - 1}] = true
	}

	// check right
	if j < (w-1) && value == (*board)[i][j+1] {
		adjacentValueCoordinates[Coordinate{i, j + 1}] = true
	}

	return adjacentValueCoordinates
}

/* TODO remove this if not needed
func deduceBlankValue(i int, j int, doneSet map[Coordinate]bool, board *[][]string) int {
	sizeOfBoundedBlank := analyzeBlankSpace(i, j, nil, doneSet, board)

	// if the size of the bounded blank is 1 or 2, we can say this point is that value
	// (if it's two, there can't be two 1's next to each other)
	// also, if we failed and can't set this, return the -1
	if sizeOfBoundedBlank == 1 || sizeOfBoundedBlank == 2 || sizeOfBoundedBlank == -1 {
		return sizeOfBoundedBlank
	} else {
		// TODO determine if we can return a value
		return -1
	}
}
*/

// determines if this blank space is part of a continuous space bound by done cells
// returns the number of blank spaces bound by done cells, or -1 if it is not bound
// TODO: targetValue is not implemented yet, the parameter is not actually used
func analyzeBlankSpace(i, j, targetValue int, alreadyAnalyzed, doneSet map[Coordinate]bool, board *[][]string) int {
	// TODO maybe don't need these as variables?  just put len(*board) or len((*board)[0]) in the if statements below
	h := len(*board)
	w := len((*board)[0])
	if alreadyAnalyzed == nil {
		alreadyAnalyzed = make(map[Coordinate]bool)
	}

	if alreadyAnalyzed[Coordinate{i, j}] {
		return 0
	}

	// first, find adjacent cells blank cells
	adjacentBlankCells := findAdjacentValues(i, j, board, "_")
	// TODO: using targetValue and maybe totalCount, use removeBadJumpingWrites to temporarily consider an
	// impossible-to-write-to cell "done" or "alreadyAnalyzed"

	// if anything is a number that is not done, then shut down the whole thing and return a -1
	// check up
	if i > 0 && !alreadyAnalyzed[Coordinate{i - 1, j}] && !adjacentBlankCells[Coordinate{i - 1, j}] && !doneSet[Coordinate{i - 1, j}] {
		return -1
	}

	// check down
	if i < (h-1) && !alreadyAnalyzed[Coordinate{i + 1, j}] && !adjacentBlankCells[Coordinate{i + 1, j}] && !doneSet[Coordinate{i + 1, j}] {
		return -1
	}

	// check left
	if j > 0 && !alreadyAnalyzed[Coordinate{i, j - 1}] && !adjacentBlankCells[Coordinate{i, j - 1}] && !doneSet[Coordinate{i, j - 1}] {
		return -1
	}

	// check right
	if j < (w-1) && !alreadyAnalyzed[Coordinate{i, j + 1}] && !adjacentBlankCells[Coordinate{i, j + 1}] && !doneSet[Coordinate{i, j + 1}] {
		return -1
	}

	// this cell can be considered bounded
	alreadyAnalyzed[Coordinate{i, j}] = true

	/*if thisCount > 100 {
		fmt.Println("overflow in blank analyzing")
	}*/

	nextCount := 0
	// if each non-blank direction is either the wall, or a done space, then recurse into any blank space
	for k := range adjacentBlankCells {
		if !alreadyAnalyzed[Coordinate{k.I, k.J}] {
			thisNextCount := analyzeBlankSpace(k.I, k.J, targetValue, alreadyAnalyzed, doneSet, board)
			if thisNextCount+nextCount > 100 {
				fmt.Println("overflow in blank analyzing")
				return -1
			}
			if thisNextCount == -1 {
				return thisNextCount
			} else {
				nextCount += thisNextCount
			}
		}
	}

	// top level doesn't have any adjacent blank bound cells
	return 1 + nextCount
}

// in a bounded blank space, analyze i, j to see if it is logically possible only to put one value, from  1 to the maxValue
// in this space, because all adjacent cells are in the doneSet and cover all other possible numbers.  So this space MUST
// be ONE of the possible values from 1 to maxValue
func determinePossibleSeedValue(i, j, maxPossibleValue int, doneSet map[Coordinate]bool, board *[][]string) int {
	// TODO maybe don't need these as variables?  just put len(*board) or len((*board)[0]) in the if statements below
	h := len(*board)
	w := len((*board)[0])

	// create and initialize presentValues map
	presentValues := make(map[int]bool)
	for x := 1; x <= maxPossibleValue; x++ {
		presentValues[x] = false
	}

	//  check up
	if i > 0 && (*board)[i-1][j] != "_" {
		if !doneSet[Coordinate{i - 1, j}] {
			return -1
		} else {
			adjacentValue, err := strconv.Atoi((*board)[i-1][j])
			if err != nil {
				fmt.Println("error")
			}
			presentValues[adjacentValue] = true
		}
	}

	// check down
	if i < (h-1) && (*board)[i+1][j] != "_" {
		if !doneSet[Coordinate{i + 1, j}] {
			return -1
		} else {
			adjacentValue, err := strconv.Atoi((*board)[i+1][j])
			if err != nil {
				fmt.Println("error")
			}
			presentValues[adjacentValue] = true
		}
	}

	// check left
	if j > 0 && (*board)[i][j-1] != "_" {
		if !doneSet[Coordinate{i, j - 1}] {
			return -1
		} else {
			adjacentValue, err := strconv.Atoi((*board)[i][j-1])
			if err != nil {
				fmt.Println("error")
			}
			presentValues[adjacentValue] = true
		}
	}

	// check right
	if j < (w-1) && (*board)[i][j+1] != "_" {
		if !doneSet[Coordinate{i, j + 1}] {
			return -1
		} else {
			adjacentValue, err := strconv.Atoi((*board)[i][j+1])
			if err != nil {
				fmt.Println("error")
			}
			presentValues[adjacentValue] = true
		}
	}

	// now see if we just have one true value in the presentValues which is unaccounted for, and if so we know this is
	// a valid seed value
	validSeedValue := -1
	falseCount := 0
	for k, v := range presentValues {
		if !v {
			falseCount++
			validSeedValue = k
		}
	}
	if falseCount == 1 {
		return validSeedValue
	} else {
		return -1
	}
}

func printBoard(b [][]string) {
	for i := 0; i < len(b); i++ {
		fmt.Printf("%s\n", strings.Join(b[i], " "))
	}
}

func main() {
	board := [][]string{
		{"8", "_", "1", "_", "_", "_", "1", "_", "_", "_"},
		{"2", "_", "_", "_", "1", "_", "2", "6", "5", "_"},
		{"_", "4", "2", "1", "_", "_", "_", "7", "_", "_"},
		{"_", "_", "_", "_", "_", "_", "2", "_", "_", "1"},
		{"_", "3", "_", "_", "4", "_", "_", "_", "1", "4"},
		{"2", "2", "_", "_", "_", "1", "_", "_", "4", "_"},
		{"1", "_", "_", "3", "_", "_", "_", "_", "_", "_"},
		{"_", "_", "3", "_", "_", "_", "1", "6", "2", "_"},
		{"_", "2", "2", "1", "_", "3", "_", "_", "_", "3"},
		{"_", "_", "_", "8", "_", "_", "_", "6", "_", "5"}}

	solvedBoard := solve(board)
	fmt.Println("No more to calculate")
	printBoard(solvedBoard)
}
