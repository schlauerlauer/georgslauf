package interfaces

import (
	"georgslauf/infra/auth"
	"georgslauf/infra/persistence"
	"georgslauf/view/host"
	"log/slog"
	"net/http"
)

type Host struct {
	repository *persistence.Repository
}

func NewHost(
	respository *persistence.Repository,
) *Host {
	return &Host{
		repository: respository,
	}
}

func (h *Host) GetHostHome() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		identity, ok := ctx.Value(auth.IdentityKey).(auth.IdentityData)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		identities, err := h.repository.Queries.GetIdentities(ctx)
		if err != nil {
			slog.Warn("could not query identities", "err", err)
		}

		w.WriteHeader(http.StatusOK)
		err = host.HostHome(host.HostHomeParams{
			Identities: identities,
			Identity:   identity,
		}).Render(ctx, w)
		if err != nil {
			slog.Warn("err rendering host home", "err", err)
		}
	})
}
