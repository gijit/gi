package shadow_gonum.org/v1/gonum/stat

import "gonum.org/v1/gonum/stat"

var Pkg = make(map[string]interface{})
func init() {
    Pkg["Bhattacharyya"] = stat.Bhattacharyya
    Pkg["BivariateMoment"] = stat.BivariateMoment
    Pkg["CDF"] = stat.CDF
    Pkg["ChiSquare"] = stat.ChiSquare
    Pkg["CircularMean"] = stat.CircularMean
    Pkg["Correlation"] = stat.Correlation
    Pkg["CorrelationMatrix"] = stat.CorrelationMatrix
    Pkg["Covariance"] = stat.Covariance
    Pkg["CovarianceMatrix"] = stat.CovarianceMatrix
    Pkg["CrossEntropy"] = stat.CrossEntropy
    Pkg["Entropy"] = stat.Entropy
    Pkg["ExKurtosis"] = stat.ExKurtosis
    Pkg["GeometricMean"] = stat.GeometricMean
    Pkg["HarmonicMean"] = stat.HarmonicMean
    Pkg["Hellinger"] = stat.Hellinger
    Pkg["Histogram"] = stat.Histogram
    Pkg["JensenShannon"] = stat.JensenShannon
    Pkg["Kendall"] = stat.Kendall
    Pkg["KolmogorovSmirnov"] = stat.KolmogorovSmirnov
    Pkg["KullbackLeibler"] = stat.KullbackLeibler
    Pkg["LinearRegression"] = stat.LinearRegression
    Pkg["Mahalanobis"] = stat.Mahalanobis
    Pkg["Mean"] = stat.Mean
    Pkg["MeanStdDev"] = stat.MeanStdDev
    Pkg["MeanVariance"] = stat.MeanVariance
    Pkg["Mode"] = stat.Mode
    Pkg["Moment"] = stat.Moment
    Pkg["MomentAbout"] = stat.MomentAbout
    Pkg["Quantile"] = stat.Quantile
    Pkg["RNoughtSquared"] = stat.RNoughtSquared
    Pkg["ROC"] = stat.ROC
    Pkg["RSquared"] = stat.RSquared
    Pkg["RSquaredFrom"] = stat.RSquaredFrom
    Pkg["Skew"] = stat.Skew
    Pkg["SortWeighted"] = stat.SortWeighted
    Pkg["SortWeightedLabeled"] = stat.SortWeightedLabeled
    Pkg["StdDev"] = stat.StdDev
    Pkg["StdErr"] = stat.StdErr
    Pkg["StdScore"] = stat.StdScore
    Pkg["Variance"] = stat.Variance

}
func GijitShadow_NewStruct_CC() *stat.CC {
	return &stat.CC{}
}


func GijitShadow_NewStruct_PC() *stat.PC {
	return &stat.PC{}
}

