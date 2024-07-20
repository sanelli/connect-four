package board

// The board
type ConnectFourBoard struct {
	content       [][]byte
	currentPlayer int
	numberOfMoves []int
	winner        int
}

func MakeConnectFourBoard() *ConnectFourBoard {
	return &ConnectFourBoard{
		make([][]byte, 6, 7), // 6 rows and 7 columns, (0,0) is bottom left corner
		1,
		make([]int, 2), // For each player the number of moves used.
		0,
	}
}

func (board *ConnectFourBoard) CurrentPlayer() int {
	return board.currentPlayer
}

func (board *ConnectFourBoard) NumberOfMoves(player int) int {
	if player < 1 || player > 2 {
		panic("Invalid player number")
	}
	return board.numberOfMoves[player-1]
}

// Get the content of the board
// Returns:
// * 0 : The cell is empty
// * 1 : The cell contains the player 1 token
// * 2 : The cell contains the player 2 token
func (board *ConnectFourBoard) GetContent(row, column int) byte {
	if column < 0 || column >= 7 {
		panic("Invalid column")
	}
	if row < 0 || row >= 6 {
		panic("Invalid row")
	}

	return board.content[row][column]
}

func (board *ConnectFourBoard) Play(column int) bool {

	if column < 0 || column >= 7 {
		panic("Invalid column")
	}

	if board.winner != 0 {
		return false
	}

	added := false
	for row := 0; row < 6; row++ {
		if board.content[row][column] == 0 {
			board.content[row][column] = byte(board.currentPlayer)
			added = true
			break
		}
	}

	if !added {
		return false
	}

	board.numberOfMoves[board.currentPlayer-1] += 1
	board.currentPlayer = (board.currentPlayer % 2) + 1
	return true
}

// Check for a winner.
// Returns:
// * -1 : The board has been filled and there is no winner
// * 0  : No player won and it is still possible playing
// * 1  : player 1 won
// * 2  : player 2 won
func (board *ConnectFourBoard) Winner() int {

	// If the winner has already been computed no need to do anything else
	if board.winner == -1 || board.winner == 1 || board.winner == 2 {
		return board.winner
	}

	// Look for horizontal
	for row := 0; row < 6; row++ {
		for column := 0; column < 4; column++ {
			for player := byte(1); player <= 2; player++ {
				if board.content[row][column] == player &&
					board.content[row][column+1] == player &&
					board.content[row][column+2] == player &&
					board.content[row][column+3] == player {
					board.winner = int(player)
					return board.winner
				}
			}
		}
	}

	// Look for vertical
	for column := 0; column < 7; column++ {
		for row := 0; row < 3; row++ {
			for player := byte(1); player <= 2; player++ {
				if board.content[row][column] == player &&
					board.content[row+1][column] == player &&
					board.content[row+2][column] == player &&
					board.content[row+3][column] == player {
					board.winner = int(player)
					return board.winner
				}
			}
		}
	}

	// Look for main digonal
	for row := 0; row < 3; row++ {
		for column := 0; column < 4; column++ {
			for player := byte(1); player <= 2; player++ {
				if board.content[row][column] == player &&
					board.content[row+1][column+1] == player &&
					board.content[row+1][column+2] == player &&
					board.content[row+1][column+3] == player {
					board.winner = int(player)
					return board.winner
				}
			}
		}
	}

	// Look for other diagonal
	for row := 3; row < 6; row++ {
		for column := 0; column < 4; column++ {
			for player := byte(1); player <= 2; player++ {
				if board.content[row][column] == player &&
					board.content[row-1][column+1] == player &&
					board.content[row-1][column+2] == player &&
					board.content[row-1][column+3] == player {
					board.winner = int(player)
					return board.winner
				}
			}
		}
	}

	// Check if the board is full
	for row := 0; row < 6; row++ {
		for column := 0; column < 7; column++ {
			if board.content[row][column] == byte(0) {
				return 0 // there is still room to play
			}
		}
	}

	board.winner = -1
	return -1
}
