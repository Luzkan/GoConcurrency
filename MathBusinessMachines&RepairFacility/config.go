package main

var BossPerf = 1
var MachPerf = 3
var WorkPerf = 4
var CustPerf = 15

var MachAddNum = 2
var MachMulNum = 2
var WorkNum = 8
var CustNum = 2
var WorkPat = 2

var MachFailure = 0.15
var RepairWorkNum = 2
var RepairWorkPerf = 3

// Error Signal & informing int set to MaxInt
var ErrorSignal = int(^uint(0) >> 1)
var NotOnList = int(^uint(0) >> 1)

var Silent = false
