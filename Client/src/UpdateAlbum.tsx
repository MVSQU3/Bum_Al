import React, { useEffect, useState } from "react";
import { useAlbumStore, type AlbumData } from "./store/useAlbumsStore";
import { useParams } from "react-router-dom";

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
    <div>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          value={formData.title}
          onChange={handleChange}
          name="title"
          placeholder="Titre"
          className="input input-accent"
        />
        <input
          type="text"
          onChange={handleChange}
          name="artist"
          value={formData.artist}
          placeholder="Artiste"
          className="input input-accent"
        />
        <input
          type="number"
          value={formData.year}
          onChange={handleChange}
          name="year"
          placeholder="Année"
          className="input input-accent"
        />
        <input
          type="file"
          onChange={handleFileChange}
          name="cover"
          accept="image/*" // Correction : "image/*" au lieu de "*/image"
          className="input"
        />
        <button type="submit" className="btn btn-active">
          Envoyé
        </button>
      </form>
    </div>
  );
};

export default UpdateAlbum;
