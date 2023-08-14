package polyseed

const (
	featureBits      = 5
	featureMask      = (1 << featureBits) - 1
	userFeatures     = 3
	userFeaturesMask = (1 << userFeatures) - 1
	encryptedMask    = 16
	reservedFeatures = featureMask ^ encryptedMask
)
