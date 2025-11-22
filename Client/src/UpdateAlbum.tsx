import React, { useEffect, useState } from "react";
import { useAlbumStore, type AlbumData } from "./store/useAlbumsStore";
import { useParams, Link } from "react-router-dom";

const UpdateAlbum = () => {
  const { id } = useParams();
  const { updateAlbum, getAlbumById } = useAlbumStore();
  const [formData, setFormData] = useState<Partial<AlbumData>>({
    title: "",
    artist: "",
    year: 2024,
    cover_url: "",
  });
  const [coverFile, setCoverFile] = useState<File | null>(null);
  const [formError, setFormError] = useState<string>("");

  useEffect(() => {
    const loadAlbum = async () => {
      let album = await getAlbumById(Number(id));
      setFormData((a) => ({ ...a, ...album }));
    };
    loadAlbum();
  }, [id]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setCoverFile(e.target.files[0]);
    }
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (!coverFile) {
      setFormError("Le fichier de couverture est requis pour mettre à jour l'album.");
      return;
    }
    setFormError("");

    // Créer un FormData pour envoyer les données multipart
    const submitData = new FormData();

    // Ajouter les champs texte
    submitData.append("title", formData.title || "");
    submitData.append("artist", formData.artist || "");
    submitData.append("year", formData.year?.toString() || "");

    // Ajouter le fichier s'il existe
    if (coverFile) {
      submitData.append("cover", coverFile);
    }

    // Appeler votre store avec FormData
    updateAlbum(submitData, Number(id));
  };
  return (
    <div className="flex items-center justify-center p-6">
      <div className="card w-full max-w-2xl bg-base-100 shadow-md">
        <div className="card-body">
          <h2 className="card-title">Mettre à jour l'album</h2>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="label">
                  <span className="label-text">Titre</span>
                </label>
                <input
                  type="text"
                  value={formData.title || ""}
                  onChange={handleChange}
                  name="title"
                  placeholder="Titre"
                  className="input input-bordered w-full"
                />

                <label className="label mt-3">
                  <span className="label-text">Artiste</span>
                </label>
                <input
                  type="text"
                  onChange={handleChange}
                  name="artist"
                  value={formData.artist || ""}
                  placeholder="Artiste"
                  className="input input-bordered w-full"
                />

                <label className="label mt-3">
                  <span className="label-text">Année</span>
                </label>
                <input
                  type="number"
                  value={formData.year ?? ""}
                  onChange={handleChange}
                  name="year"
                  placeholder="Année"
                  className="input input-bordered w-full"
                />
              </div>

              <div>
                <label className="label">
                  <span className="label-text">Couverture</span>
                </label>

                {formData.cover_url ? (
                  <div className="mb-3">
                    <img
                      src={formData.cover_url}
                      alt={formData.title || "cover"}
                      className="w-full h-48 object-cover rounded"
                    />
                  </div>
                ) : (
                  <div className="mb-3 h-48 w-full bg-base-200 flex items-center justify-center rounded">
                    <span className="text-sm text-muted">Aucune image</span>
                  </div>
                )}

                <input
                  type="file"
                  onChange={handleFileChange}
                  name="cover"
                  accept="image/*"
                  className="file-input file-input-bordered w-full"
                />
                {!coverFile && (
                  <p className="text-sm text-error mt-2">Veuillez sélectionner un fichier de couverture pour pouvoir enregistrer.</p>
                )}
              </div>
            </div>

            <div className="flex items-center justify-end gap-3">
              <Link to="/" className="btn btn-ghost">
                Annuler
              </Link>
              <button type="submit" className="btn btn-primary" disabled={!coverFile}>
                Enregistrer
              </button>
            </div>
          </form>
          {formError && <p className="text-sm text-error mt-3">{formError}</p>}
        </div>
      </div>
    </div>
  );
};

export default UpdateAlbum;
