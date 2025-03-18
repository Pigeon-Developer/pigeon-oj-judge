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

func (source SourceDB) GetOne(languageList []int) (*Solution, error) {
	var err error
	solution := SolutionRecord{}

	// SELECT solution_id FROM solution WHERE language in (%s) and result<2 ORDER BY result, solution_id  limit %d
	query, args, err := sqlx.In("SELECT * FROM solution WHERE result < 2 AND language IN (?) ORDER BY result, solution_id limit 1", languageList)
	if err != nil {
		return nil, err
	}
	err = source.db.Get(&solution, query, args...)
	if err != nil {
		return nil, err
	}

	source.Update(solution.SolutionId, SolutionResult{Result: Result_CI})

	problem := ProblemRecord{}

	err = source.db.Get(&problem, "SELECT * FROM problem WHERE problem_id = ? LIMIT 1", solution.ProblemId)
	if err != nil {
		return nil, err
	}

	ret := Solution{
		SolutionId:  solution.SolutionId,
		ProblemId:   solution.ProblemId,
		ContestId:   solution.ContestId,
		UserId:      solution.UserId,
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

const JudgerName = "pigeon-oj-judge"

func (source SourceDB) Update(solutionId int, result SolutionResult) error {
	switch result.Result {
	case Result_CI, Result_RJ:
		{
			// 这部分仅更新
			updateSolutionSql := `UPDATE solution SET result = ? , judger = ? WHERE solution_id = ?`
			_, err := source.db.Exec(updateSolutionSql, result.Result, "pigeon-oj-judge", solutionId)

			return err
		}
	}

	// 更新判题开销
	updateSolutionSql := `UPDATE solution SET result = ? , time= ?, memory= ?, judger = ?, judgetime = now() WHERE solution_id = ?`
	_, err := source.db.Exec(updateSolutionSql, result.Result, result.TimeCost, result.MemoryUsage, "pigeon-oj-judge", solutionId)

	// 更新错误信息
	switch result.Result {
	case Result_CE:
		// compileinfo
		deleteInfoSql := `DELETE FROM compileinfo WHERE solution_id = ?`
		source.db.Exec(deleteInfoSql, solutionId)

		appendInfoSql := `INSERT INTO compileinfo VALUES(? , ?)`
		source.db.Exec(appendInfoSql, solutionId, result.Info)
	case Result_RE:
		// runtimeinfo
		deleteInfoSql := `DELETE FROM runtimeinfo WHERE solution_id = ?`
		source.db.Exec(deleteInfoSql, solutionId)

		appendInfoSql := `INSERT INTO runtimeinfo VALUES(? , ?)`
		source.db.Exec(appendInfoSql, solutionId, result.Info)
	}

	// 更新用户数据
	if result.Result == Result_AC {
		updateSolved := "UPDATE `users` SET `solved`=(SELECT count(DISTINCT `problem_id`) FROM `solution` WHERE `user_id`= ? AND `result`=4) WHERE `user_id`= ?"
		source.db.Exec(updateSolved, result.Solution.UserId, result.Solution.UserId)

		updateSubmit := "UPDATE `users` SET `submit`=(SELECT count(1) FROM `solution` WHERE `user_id`= ? and problem_id > 0) WHERE `user_id`= ?"
		source.db.Exec(updateSubmit, result.Solution.UserId, result.Solution.UserId)
	}

	// 更新题目数据
	{
		updateAccepted := "UPDATE `problem` SET `accepted`=(SELECT count(1) FROM `solution` WHERE `problem_id`= ? AND `result`=4) WHERE `problem_id`= ?"
		source.db.Exec(updateAccepted, result.Solution.ProblemId, result.Solution.ProblemId)
	}

	// @TODO 处理比赛的情况
	return err
}

func (source SourceDB) Close() {

}
