package constants

import "errors"

var NotImplemented error = errors.New("not implemented")
var InvalidParam error = errors.New("invalid param")
var TimeOut error = errors.New("invalid param")
var InvalidType error = errors.New("not implemented")
var InvalidSignature error = errors.New("InvalidSignature")
var InvalidState error = errors.New("InvalidState")
var InvalidSequenceNumber error = errors.New("InvalidSequenceNumber")
var DuplicateKey error = errors.New("DuplicateKey")

const BottomStackName = "BottomStack"
const TopStackName = "TopStack"
