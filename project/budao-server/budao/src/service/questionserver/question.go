package questionserver

import (
	"context"
	"fmt"
	"service/util"
	"strconv"

	"github.com/sumaig/glog"

	"common"
	"db"
	pb "twirprpc"
)

// Server identify for QuestionService RPC
type Server struct{}

// GetServer return server of question service
func GetServer() *Server {
	server := &Server{}

	return server
}

// AnswerQuestion report user answer
func (s *Server) AnswerQuestion(ctx context.Context, req *pb.AnswerQuestionRequest) (resp *pb.AnswerQuestionResponse, err error) {
	glog.Debug("req:%v", req)
	resp = &pb.AnswerQuestionResponse{
		Status: common.GetInitStatus(),
	}
	defer func() {
		if err != nil {
			resp.Status.Message = err.Error()
			err = nil
		}
		glog.Debug("resp:%v", resp)
	}()

	userID, _ := strconv.ParseUint(req.GetHeader().GetUserId(), 10, 64)
	questionID, _ := strconv.ParseUint(req.GetQuestionId(), 10, 64)
	optionID, _ := strconv.ParseUint(req.GetOptionId(), 10, 64)
	token := req.GetHeader().GetToken()

	var valid bool
	if valid, err = util.CheckIDValid(userID, token); !valid {
		resp.Status.Code = pb.Status_USER_NOTLOGIN
		glog.Error("userid:%d, token:%s check not valid. err:%v", userID, token, err)
		return
	}

	// 1. Check if the question has been answered
	key := fmt.Sprintf("%s%s", common.USER_ACT_PREFIX, strconv.FormatUint(userID, 10))
	ret, err := db.HGet(key, strconv.FormatUint(questionID, 10))
	if ret != "" {
		sqlString := fmt.Sprintf("select option_content, is_answer, answer_num from question_option where question_id=%d and option_id=%s", questionID, ret)
		var content string
		var isAnswer uint32
		var answerNum uint32
		var right bool
		tempRow, _ := db.QueryRow(common.BUDAODB, sqlString)
		err = tempRow.Scan(&content, &isAnswer, &answerNum)
		if isAnswer == 1 {
			right = true
		}
		optionItem := &pb.OptionItem{
			OptionId:    ret,
			QuestionId:  strconv.FormatUint(questionID, 10),
			Content:     content,
			Right:       right,
			ChooseCount: answerNum,
		}
		resp.Option = optionItem
		resp.Status.Message = "question has been answered"
		resp.Status.Code = pb.Status_OK
		return
	}

	// 2. get the question answer
	sqlString := fmt.Sprintf("select is_answer from question_option where question_id = %d and option_id = %d", questionID, optionID)
	var answerFlag int
	firstRow, err := db.QueryRow(common.BUDAODB, sqlString)
	err = firstRow.Scan(&answerFlag)
	if err != nil {
		glog.Error("query question answer failed. err:%v", err)
		return
	}

	// 3. get the question score
	sqlString = fmt.Sprintf("select score, vid from question where id = %d", questionID)
	var score uint32
	var vid uint64
	secondRow, err := db.QueryRow(common.BUDAODB, sqlString)
	err = secondRow.Scan(&score, &vid)
	if err != nil {
		glog.Error("query question score failed. err:%v", err)
		return
	}

	var (
		userResult int
		userScore  uint32
	)
	userTableName, err := db.GetTableName("user_", userID)
	tableNumber := userID >> 54
	userQuestionTN := fmt.Sprintf("user_question_%d", tableNumber)
	if answerFlag == 1 {
		// answer correct
		sqlString = fmt.Sprintf("update %s set right_answer_num = right_answer_num+1, get_score = get_score+%d where uid = %d", userTableName, score, userID)
		_, err = db.Exec(common.BUDAODB, sqlString)
		if err != nil {
			glog.Error("update user table failed. err:%v", err)
			return
		}

		sqlString = fmt.Sprintf("update question set right_answer_num = right_answer_num+1 where id = %d", questionID)
		_, err = db.Exec(common.BUDAODB, sqlString)
		if err != nil {
			glog.Error("update question table failed. err:%v", err)
			return
		}

		userScore = score
		userResult = 1
	} else {
		// answer error
		sqlString = fmt.Sprintf("update %s set wrong_answer_num = wrong_answer_num+1 where uid = %d", userTableName, userID)
		_, err = db.Exec(common.BUDAODB, sqlString)
		if err != nil {
			glog.Error("update user table failed. err:%v", err)
			return
		}

		sqlString = fmt.Sprintf("update question set wrong_answer_num = wrong_answer_num+1 where id = %d", questionID)
		_, err = db.Exec(common.BUDAODB, sqlString)
		if err != nil {
			glog.Error("update question table failed. err:%v", err)
			return
		}

		userScore = 0
		userResult = 2
	}

	sqlString = fmt.Sprintf("update question_option set answer_num = answer_num+1 where question_id = %d and option_id = %d", questionID, optionID)
	_, err = db.Exec(common.BUDAODB, sqlString)
	if err != nil {
		glog.Error("update question_option table failed. err:%v", err)
		return
	}

	sqlString = fmt.Sprintf("insert into %s (uid, question_id, result, get_score, option_id) values (%d, %d, %d, %d, %d)", userQuestionTN, userID, questionID, userResult, userScore, optionID)
	_, err = db.Exec(common.BUDAODB, sqlString)
	if err != nil {
		glog.Error("answer question insert user_question_ table faild. err:%v", err)
		return
	}

	// update question_dynamic
	field := fmt.Sprintf("%s%s", strconv.FormatUint(questionID, 10), strconv.FormatUint(optionID, 10))
	_, err = db.HIncrBy(common.QUESTION_DYNAMIC, field, 1)
	if err != nil {
		glog.Error("answer question insert hash question_dynamic faild. err:%v", err)
	}

	// inster user_act_[uid]
	_, err = db.HSet(key, strconv.FormatUint(questionID, 10), strconv.FormatUint(optionID, 10))
	if err != nil {
		glog.Error("answer question insert hash user_act_ faild. err:%v", err)
	}

	// update vid_dynamic
	_, err = db.SAdd(common.VID_DYNAMIC, strconv.FormatUint(vid, 10))
	if err != nil {
		glog.Error("insert vid_dynamic set failed. err:%v", err)
		return
	}

	resp.Status.Code = pb.Status_OK

	return
}
