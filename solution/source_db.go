package solution

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type SourceDB struct {
	dsn    string
	dbType string
	db     *sqlx.DB
}

// 从 db 直接获取未判题的提交
func NewSolutionSourceDB(dbType, dsn string) (*SourceDB, error) {
	source := SourceDB{
		dsn:    dsn,
		dbType: dbType,
	}

	db, err := sqlx.Connect(source.dbType, source.dsn)

	if err != nil {
		return nil, err
	}

	source.db = db.Unsafe()

	return &source, nil
}

func (source SourceDB) GetOne() (*Solution, error) {
	var err error
	solution := SolutionRecord{}

	// SELECT solution_id FROM solution WHERE language in (%s) and result<2 ORDER BY result, solution_id  limit %d
	err = source.db.Get(&solution, "SELECT * FROM solution WHERE result < 2 ORDER BY result, solution_id  limit 1")
	if err != nil {
		return nil, err
	}

	problem := ProblemRecord{}

	err = source.db.Get(&problem, "SELECT * FROM problem WHERE problem_id = ? LIMIT 1", solution.ProblemId)
	if err != nil {
		return nil, err
	}

	ret := Solution{
		SolutionId:  solution.SolutionId,
		ProblemId:   solution.ProblemId,
		ContestId:   solution.ContestId,
		Language:    solution.Language,
		TimeLimit:   problem.TimeLimit,
		MemoryLimit: problem.MemoryLimit,
		Code:        "",
	}

	sourceCode := SourceCodeRecord{}

	err = source.db.Get(&sourceCode, "SELECT * FROM source_code WHERE solution_id = ? LIMIT 1", solution.SolutionId)

	if err != nil {
		return nil, err
	}

	ret.Code = sourceCode.Source

	return &ret, nil
}

func (source SourceDB) Update(solutionId int, result SolutionResult) error {
	_, err := source.db.Exec(`UPDATE solution SET result = ? WHERE solution_id = ?`, result.Result, solutionId)
	return err
}

func (source SourceDB) Close() {

}
