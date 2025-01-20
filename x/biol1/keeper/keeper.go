package keeper

type Keeper struct {
	sidecarURL string
}

func NewKeeper(sidecarURL string) Keeper {
	return Keeper{sidecarURL: sidecarURL}
}

func (k Keeper) SidecarURL() string {
	return k.sidecarURL
}
