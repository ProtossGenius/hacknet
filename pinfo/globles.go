package pinfo

// this file to store global vars, for show less error when coding.

// pointInfoMgrFactory a public pointInfoMgrFactory to create PointInfoMgrItf.
var pointInfoMgrFactory PointInfoMgrFactory = newMemoryMgr

// p2PHelperFacotry create p2p helper.
var p2PHelperFacotry P2PHelperFacotry = newP2pHelperMemoryImpl
