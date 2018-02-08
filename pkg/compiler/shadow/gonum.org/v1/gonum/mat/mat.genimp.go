package shadow_gonum.org/v1/gonum/mat

import "gonum.org/v1/gonum/mat"

var Pkg = make(map[string]interface{})
func init() {
    Pkg["BandWidther"] = GijitShadow_InterfaceConvertTo2_BandWidther
    Pkg["Banded"] = GijitShadow_InterfaceConvertTo2_Banded
    Pkg["CMatrix"] = GijitShadow_InterfaceConvertTo2_CMatrix
    Pkg["Cloner"] = GijitShadow_InterfaceConvertTo2_Cloner
    Pkg["Col"] = mat.Col
    Pkg["ColNonZeroDoer"] = GijitShadow_InterfaceConvertTo2_ColNonZeroDoer
    Pkg["ColViewer"] = GijitShadow_InterfaceConvertTo2_ColViewer
    Pkg["Cond"] = mat.Cond
    Pkg["ConditionTolerance"] = mat.ConditionTolerance
    Pkg["Copier"] = GijitShadow_InterfaceConvertTo2_Copier
    Pkg["DenseCopyOf"] = mat.DenseCopyOf
    Pkg["Det"] = mat.Det
    Pkg["Dot"] = mat.Dot
    Pkg["DotByte"] = mat.DotByte
    Pkg["Equal"] = mat.Equal
    Pkg["EqualApprox"] = mat.EqualApprox
    Pkg["ErrBandSet"] = mat.ErrBandSet
    Pkg["ErrColAccess"] = mat.ErrColAccess
    Pkg["ErrColLength"] = mat.ErrColLength
    Pkg["ErrFailedEigen"] = mat.ErrFailedEigen
    Pkg["ErrIllegalStride"] = mat.ErrIllegalStride
    Pkg["ErrIndexOutOfRange"] = mat.ErrIndexOutOfRange
    Pkg["ErrNormOrder"] = mat.ErrNormOrder
    Pkg["ErrNotPSD"] = mat.ErrNotPSD
    Pkg["ErrPivot"] = mat.ErrPivot
    Pkg["ErrRowAccess"] = mat.ErrRowAccess
    Pkg["ErrRowLength"] = mat.ErrRowLength
    Pkg["ErrShape"] = mat.ErrShape
    Pkg["ErrSingular"] = mat.ErrSingular
    Pkg["ErrSliceLengthMismatch"] = mat.ErrSliceLengthMismatch
    Pkg["ErrSquare"] = mat.ErrSquare
    Pkg["ErrTriangle"] = mat.ErrTriangle
    Pkg["ErrTriangleSet"] = mat.ErrTriangleSet
    Pkg["ErrVectorAccess"] = mat.ErrVectorAccess
    Pkg["ErrZeroLength"] = mat.ErrZeroLength
    Pkg["Excerpt"] = mat.Excerpt
    Pkg["Formatted"] = mat.Formatted
    Pkg["Grower"] = GijitShadow_InterfaceConvertTo2_Grower
    Pkg["Inner"] = mat.Inner
    Pkg["LogDet"] = mat.LogDet
    Pkg["Matrix"] = GijitShadow_InterfaceConvertTo2_Matrix
    Pkg["Max"] = mat.Max
    Pkg["Maybe"] = mat.Maybe
    Pkg["MaybeComplex"] = mat.MaybeComplex
    Pkg["MaybeFloat"] = mat.MaybeFloat
    Pkg["Min"] = mat.Min
    Pkg["Mutable"] = GijitShadow_InterfaceConvertTo2_Mutable
    Pkg["MutableBanded"] = GijitShadow_InterfaceConvertTo2_MutableBanded
    Pkg["MutableSymBanded"] = GijitShadow_InterfaceConvertTo2_MutableSymBanded
    Pkg["MutableSymmetric"] = GijitShadow_InterfaceConvertTo2_MutableSymmetric
    Pkg["MutableTriangular"] = GijitShadow_InterfaceConvertTo2_MutableTriangular
    Pkg["NewBandDense"] = mat.NewBandDense
    Pkg["NewDense"] = mat.NewDense
    Pkg["NewDiagonal"] = mat.NewDiagonal
    Pkg["NewDiagonalRect"] = mat.NewDiagonalRect
    Pkg["NewSymBandDense"] = mat.NewSymBandDense
    Pkg["NewSymDense"] = mat.NewSymDense
    Pkg["NewTriDense"] = mat.NewTriDense
    Pkg["NewVecDense"] = mat.NewVecDense
    Pkg["NonZeroDoer"] = GijitShadow_InterfaceConvertTo2_NonZeroDoer
    Pkg["Norm"] = mat.Norm
    Pkg["Prefix"] = mat.Prefix
    Pkg["RawBander"] = GijitShadow_InterfaceConvertTo2_RawBander
    Pkg["RawColViewer"] = GijitShadow_InterfaceConvertTo2_RawColViewer
    Pkg["RawMatrixSetter"] = GijitShadow_InterfaceConvertTo2_RawMatrixSetter
    Pkg["RawMatrixer"] = GijitShadow_InterfaceConvertTo2_RawMatrixer
    Pkg["RawRowViewer"] = GijitShadow_InterfaceConvertTo2_RawRowViewer
    Pkg["RawSymBander"] = GijitShadow_InterfaceConvertTo2_RawSymBander
    Pkg["RawSymmetricer"] = GijitShadow_InterfaceConvertTo2_RawSymmetricer
    Pkg["RawTriangular"] = GijitShadow_InterfaceConvertTo2_RawTriangular
    Pkg["RawVectorer"] = GijitShadow_InterfaceConvertTo2_RawVectorer
    Pkg["Reseter"] = GijitShadow_InterfaceConvertTo2_Reseter
    Pkg["Row"] = mat.Row
    Pkg["RowNonZeroDoer"] = GijitShadow_InterfaceConvertTo2_RowNonZeroDoer
    Pkg["RowViewer"] = GijitShadow_InterfaceConvertTo2_RowViewer
    Pkg["Squeeze"] = mat.Squeeze
    Pkg["Sum"] = mat.Sum
    Pkg["Symmetric"] = GijitShadow_InterfaceConvertTo2_Symmetric
    Pkg["Trace"] = mat.Trace
    Pkg["Triangular"] = GijitShadow_InterfaceConvertTo2_Triangular
    Pkg["Unconjugator"] = GijitShadow_InterfaceConvertTo2_Unconjugator
    Pkg["UntransposeBander"] = GijitShadow_InterfaceConvertTo2_UntransposeBander
    Pkg["UntransposeTrier"] = GijitShadow_InterfaceConvertTo2_UntransposeTrier
    Pkg["Untransposer"] = GijitShadow_InterfaceConvertTo2_Untransposer
    Pkg["VecDenseCopyOf"] = mat.VecDenseCopyOf
    Pkg["Vector"] = GijitShadow_InterfaceConvertTo2_Vector

}
func GijitShadow_NewStruct_BandDense() *mat.BandDense {
	return &mat.BandDense{}
}


func GijitShadow_InterfaceConvertTo2_BandWidther(x interface{}) (y mat.BandWidther, b bool) {
	y, b = x.(mat.BandWidther)
	return
}

func GijitShadow_InterfaceConvertTo1_BandWidther(x interface{}) mat.BandWidther {
	return x.(mat.BandWidther)
}


func GijitShadow_InterfaceConvertTo2_Banded(x interface{}) (y mat.Banded, b bool) {
	y, b = x.(mat.Banded)
	return
}

func GijitShadow_InterfaceConvertTo1_Banded(x interface{}) mat.Banded {
	return x.(mat.Banded)
}


func GijitShadow_InterfaceConvertTo2_CMatrix(x interface{}) (y mat.CMatrix, b bool) {
	y, b = x.(mat.CMatrix)
	return
}

func GijitShadow_InterfaceConvertTo1_CMatrix(x interface{}) mat.CMatrix {
	return x.(mat.CMatrix)
}


func GijitShadow_NewStruct_Cholesky() *mat.Cholesky {
	return &mat.Cholesky{}
}


func GijitShadow_InterfaceConvertTo2_Cloner(x interface{}) (y mat.Cloner, b bool) {
	y, b = x.(mat.Cloner)
	return
}

func GijitShadow_InterfaceConvertTo1_Cloner(x interface{}) mat.Cloner {
	return x.(mat.Cloner)
}


func GijitShadow_InterfaceConvertTo2_ColNonZeroDoer(x interface{}) (y mat.ColNonZeroDoer, b bool) {
	y, b = x.(mat.ColNonZeroDoer)
	return
}

func GijitShadow_InterfaceConvertTo1_ColNonZeroDoer(x interface{}) mat.ColNonZeroDoer {
	return x.(mat.ColNonZeroDoer)
}


func GijitShadow_InterfaceConvertTo2_ColViewer(x interface{}) (y mat.ColViewer, b bool) {
	y, b = x.(mat.ColViewer)
	return
}

func GijitShadow_InterfaceConvertTo1_ColViewer(x interface{}) mat.ColViewer {
	return x.(mat.ColViewer)
}


func GijitShadow_NewStruct_Conjugate() *mat.Conjugate {
	return &mat.Conjugate{}
}


func GijitShadow_InterfaceConvertTo2_Copier(x interface{}) (y mat.Copier, b bool) {
	y, b = x.(mat.Copier)
	return
}

func GijitShadow_InterfaceConvertTo1_Copier(x interface{}) mat.Copier {
	return x.(mat.Copier)
}


func GijitShadow_NewStruct_Dense() *mat.Dense {
	return &mat.Dense{}
}


func GijitShadow_NewStruct_Eigen() *mat.Eigen {
	return &mat.Eigen{}
}


func GijitShadow_NewStruct_EigenSym() *mat.EigenSym {
	return &mat.EigenSym{}
}


func GijitShadow_NewStruct_Error() *mat.Error {
	return &mat.Error{}
}


func GijitShadow_NewStruct_ErrorStack() *mat.ErrorStack {
	return &mat.ErrorStack{}
}


func GijitShadow_NewStruct_GSVD() *mat.GSVD {
	return &mat.GSVD{}
}


func GijitShadow_InterfaceConvertTo2_Grower(x interface{}) (y mat.Grower, b bool) {
	y, b = x.(mat.Grower)
	return
}

func GijitShadow_InterfaceConvertTo1_Grower(x interface{}) mat.Grower {
	return x.(mat.Grower)
}


func GijitShadow_NewStruct_HOGSVD() *mat.HOGSVD {
	return &mat.HOGSVD{}
}


func GijitShadow_NewStruct_LQ() *mat.LQ {
	return &mat.LQ{}
}


func GijitShadow_NewStruct_LU() *mat.LU {
	return &mat.LU{}
}


func GijitShadow_InterfaceConvertTo2_Matrix(x interface{}) (y mat.Matrix, b bool) {
	y, b = x.(mat.Matrix)
	return
}

func GijitShadow_InterfaceConvertTo1_Matrix(x interface{}) mat.Matrix {
	return x.(mat.Matrix)
}


func GijitShadow_InterfaceConvertTo2_Mutable(x interface{}) (y mat.Mutable, b bool) {
	y, b = x.(mat.Mutable)
	return
}

func GijitShadow_InterfaceConvertTo1_Mutable(x interface{}) mat.Mutable {
	return x.(mat.Mutable)
}


func GijitShadow_InterfaceConvertTo2_MutableBanded(x interface{}) (y mat.MutableBanded, b bool) {
	y, b = x.(mat.MutableBanded)
	return
}

func GijitShadow_InterfaceConvertTo1_MutableBanded(x interface{}) mat.MutableBanded {
	return x.(mat.MutableBanded)
}


func GijitShadow_InterfaceConvertTo2_MutableSymBanded(x interface{}) (y mat.MutableSymBanded, b bool) {
	y, b = x.(mat.MutableSymBanded)
	return
}

func GijitShadow_InterfaceConvertTo1_MutableSymBanded(x interface{}) mat.MutableSymBanded {
	return x.(mat.MutableSymBanded)
}


func GijitShadow_InterfaceConvertTo2_MutableSymmetric(x interface{}) (y mat.MutableSymmetric, b bool) {
	y, b = x.(mat.MutableSymmetric)
	return
}

func GijitShadow_InterfaceConvertTo1_MutableSymmetric(x interface{}) mat.MutableSymmetric {
	return x.(mat.MutableSymmetric)
}


func GijitShadow_InterfaceConvertTo2_MutableTriangular(x interface{}) (y mat.MutableTriangular, b bool) {
	y, b = x.(mat.MutableTriangular)
	return
}

func GijitShadow_InterfaceConvertTo1_MutableTriangular(x interface{}) mat.MutableTriangular {
	return x.(mat.MutableTriangular)
}


func GijitShadow_InterfaceConvertTo2_NonZeroDoer(x interface{}) (y mat.NonZeroDoer, b bool) {
	y, b = x.(mat.NonZeroDoer)
	return
}

func GijitShadow_InterfaceConvertTo1_NonZeroDoer(x interface{}) mat.NonZeroDoer {
	return x.(mat.NonZeroDoer)
}


func GijitShadow_NewStruct_QR() *mat.QR {
	return &mat.QR{}
}


func GijitShadow_InterfaceConvertTo2_RawBander(x interface{}) (y mat.RawBander, b bool) {
	y, b = x.(mat.RawBander)
	return
}

func GijitShadow_InterfaceConvertTo1_RawBander(x interface{}) mat.RawBander {
	return x.(mat.RawBander)
}


func GijitShadow_InterfaceConvertTo2_RawColViewer(x interface{}) (y mat.RawColViewer, b bool) {
	y, b = x.(mat.RawColViewer)
	return
}

func GijitShadow_InterfaceConvertTo1_RawColViewer(x interface{}) mat.RawColViewer {
	return x.(mat.RawColViewer)
}


func GijitShadow_InterfaceConvertTo2_RawMatrixSetter(x interface{}) (y mat.RawMatrixSetter, b bool) {
	y, b = x.(mat.RawMatrixSetter)
	return
}

func GijitShadow_InterfaceConvertTo1_RawMatrixSetter(x interface{}) mat.RawMatrixSetter {
	return x.(mat.RawMatrixSetter)
}


func GijitShadow_InterfaceConvertTo2_RawMatrixer(x interface{}) (y mat.RawMatrixer, b bool) {
	y, b = x.(mat.RawMatrixer)
	return
}

func GijitShadow_InterfaceConvertTo1_RawMatrixer(x interface{}) mat.RawMatrixer {
	return x.(mat.RawMatrixer)
}


func GijitShadow_InterfaceConvertTo2_RawRowViewer(x interface{}) (y mat.RawRowViewer, b bool) {
	y, b = x.(mat.RawRowViewer)
	return
}

func GijitShadow_InterfaceConvertTo1_RawRowViewer(x interface{}) mat.RawRowViewer {
	return x.(mat.RawRowViewer)
}


func GijitShadow_InterfaceConvertTo2_RawSymBander(x interface{}) (y mat.RawSymBander, b bool) {
	y, b = x.(mat.RawSymBander)
	return
}

func GijitShadow_InterfaceConvertTo1_RawSymBander(x interface{}) mat.RawSymBander {
	return x.(mat.RawSymBander)
}


func GijitShadow_InterfaceConvertTo2_RawSymmetricer(x interface{}) (y mat.RawSymmetricer, b bool) {
	y, b = x.(mat.RawSymmetricer)
	return
}

func GijitShadow_InterfaceConvertTo1_RawSymmetricer(x interface{}) mat.RawSymmetricer {
	return x.(mat.RawSymmetricer)
}


func GijitShadow_InterfaceConvertTo2_RawTriangular(x interface{}) (y mat.RawTriangular, b bool) {
	y, b = x.(mat.RawTriangular)
	return
}

func GijitShadow_InterfaceConvertTo1_RawTriangular(x interface{}) mat.RawTriangular {
	return x.(mat.RawTriangular)
}


func GijitShadow_InterfaceConvertTo2_RawVectorer(x interface{}) (y mat.RawVectorer, b bool) {
	y, b = x.(mat.RawVectorer)
	return
}

func GijitShadow_InterfaceConvertTo1_RawVectorer(x interface{}) mat.RawVectorer {
	return x.(mat.RawVectorer)
}


func GijitShadow_InterfaceConvertTo2_Reseter(x interface{}) (y mat.Reseter, b bool) {
	y, b = x.(mat.Reseter)
	return
}

func GijitShadow_InterfaceConvertTo1_Reseter(x interface{}) mat.Reseter {
	return x.(mat.Reseter)
}


func GijitShadow_InterfaceConvertTo2_RowNonZeroDoer(x interface{}) (y mat.RowNonZeroDoer, b bool) {
	y, b = x.(mat.RowNonZeroDoer)
	return
}

func GijitShadow_InterfaceConvertTo1_RowNonZeroDoer(x interface{}) mat.RowNonZeroDoer {
	return x.(mat.RowNonZeroDoer)
}


func GijitShadow_InterfaceConvertTo2_RowViewer(x interface{}) (y mat.RowViewer, b bool) {
	y, b = x.(mat.RowViewer)
	return
}

func GijitShadow_InterfaceConvertTo1_RowViewer(x interface{}) mat.RowViewer {
	return x.(mat.RowViewer)
}


func GijitShadow_NewStruct_SVD() *mat.SVD {
	return &mat.SVD{}
}


func GijitShadow_NewStruct_SymBandDense() *mat.SymBandDense {
	return &mat.SymBandDense{}
}


func GijitShadow_NewStruct_SymDense() *mat.SymDense {
	return &mat.SymDense{}
}


func GijitShadow_InterfaceConvertTo2_Symmetric(x interface{}) (y mat.Symmetric, b bool) {
	y, b = x.(mat.Symmetric)
	return
}

func GijitShadow_InterfaceConvertTo1_Symmetric(x interface{}) mat.Symmetric {
	return x.(mat.Symmetric)
}


func GijitShadow_NewStruct_Transpose() *mat.Transpose {
	return &mat.Transpose{}
}


func GijitShadow_NewStruct_TransposeBand() *mat.TransposeBand {
	return &mat.TransposeBand{}
}


func GijitShadow_NewStruct_TransposeTri() *mat.TransposeTri {
	return &mat.TransposeTri{}
}


func GijitShadow_NewStruct_TransposeVec() *mat.TransposeVec {
	return &mat.TransposeVec{}
}


func GijitShadow_NewStruct_TriDense() *mat.TriDense {
	return &mat.TriDense{}
}


func GijitShadow_InterfaceConvertTo2_Triangular(x interface{}) (y mat.Triangular, b bool) {
	y, b = x.(mat.Triangular)
	return
}

func GijitShadow_InterfaceConvertTo1_Triangular(x interface{}) mat.Triangular {
	return x.(mat.Triangular)
}


func GijitShadow_InterfaceConvertTo2_Unconjugator(x interface{}) (y mat.Unconjugator, b bool) {
	y, b = x.(mat.Unconjugator)
	return
}

func GijitShadow_InterfaceConvertTo1_Unconjugator(x interface{}) mat.Unconjugator {
	return x.(mat.Unconjugator)
}


func GijitShadow_InterfaceConvertTo2_UntransposeBander(x interface{}) (y mat.UntransposeBander, b bool) {
	y, b = x.(mat.UntransposeBander)
	return
}

func GijitShadow_InterfaceConvertTo1_UntransposeBander(x interface{}) mat.UntransposeBander {
	return x.(mat.UntransposeBander)
}


func GijitShadow_InterfaceConvertTo2_UntransposeTrier(x interface{}) (y mat.UntransposeTrier, b bool) {
	y, b = x.(mat.UntransposeTrier)
	return
}

func GijitShadow_InterfaceConvertTo1_UntransposeTrier(x interface{}) mat.UntransposeTrier {
	return x.(mat.UntransposeTrier)
}


func GijitShadow_InterfaceConvertTo2_Untransposer(x interface{}) (y mat.Untransposer, b bool) {
	y, b = x.(mat.Untransposer)
	return
}

func GijitShadow_InterfaceConvertTo1_Untransposer(x interface{}) mat.Untransposer {
	return x.(mat.Untransposer)
}


func GijitShadow_NewStruct_VecDense() *mat.VecDense {
	return &mat.VecDense{}
}


func GijitShadow_InterfaceConvertTo2_Vector(x interface{}) (y mat.Vector, b bool) {
	y, b = x.(mat.Vector)
	return
}

func GijitShadow_InterfaceConvertTo1_Vector(x interface{}) mat.Vector {
	return x.(mat.Vector)
}

