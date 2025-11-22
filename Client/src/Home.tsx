import { useEffect } from "react";
import { useAlbumStore } from "./store/useAlbumsStore";
import { Link } from "react-router-dom";
import { Loader, LogOut } from "lucide-react";
import { useAuthStore } from "./store/userAuthStore";

const Home = () => {
  const { getAllAlbums, deleteAlbum, albums, isLoading } = useAlbumStore();
  const { logout } = useAuthStore();

  const handleDelete = () => {
    if (window.confirm("Are your sure ?")) {
      logout();
    }
  };
  useEffect(() => {
    getAllAlbums();
  }, []);
  return (
    <div className="p-4">
      <div className="flex items-center justify-between mb-4">
        <form>
          <input
            type="text"
            name=""
            placeholder="Search"
            className="input input-neutral"
          />
        </form>

        <div className="flex items-center gap-2">
          <Link to="/album/add" className="btn btn-primary">
            Ajouter un album
          </Link>
          <button
            type="button"
            onClick={handleDelete}
            title="Rafraîchir"
            className="btn btn-square btn-sm btn-outline btn-error"
          >
            <LogOut />
          </button>
        </div>
      </div>

      {isLoading ? (
        <p>
          <Loader />{" "}
        </p>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {albums && albums.length > 0 ? (
            albums.map((a) => (
              <div key={a.id} className="border rounded p-3 shadow-sm">
                {a.cover_url && (
                  <img
                    src={a.cover_url}
                    alt={a.title}
                    className="w-full h-48 object-cover mb-2"
                  />
                )}
                <h3 className="font-semibold">{a.title}</h3>
                <p className="text-sm text-muted">
                  {a.artist} — {a.year}
                </p>
                <div className="mt-3 flex justify-between gap-1">
                  <Link
                    to={`/album/update/${a.id}`}
                    className="btn btn-secondary btn-sm"
                  >
                    Mettre à jour
                  </Link>
                  <button
                    type="button"
                    className="btn btn-error btn-sm"
                    onClick={() => {
                      const ok = window.confirm(
                        `Supprimer l'album "${a.title}" ? Cette action est irréversible.`
                      );
                      if (ok) deleteAlbum(a.id);
                    }}
                  >
                    Supprimer
                  </button>
                </div>
              </div>
            ))
          ) : (
            <p>Aucun album trouvé.</p>
          )}
        </div>
      )}
    </div>
  );
};

export default Home;
