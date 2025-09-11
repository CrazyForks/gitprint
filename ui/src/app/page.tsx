"use client";

import { useEffect, useState } from "react";
import Form from "../components/Form";
import { getRecentRepos } from "../lib/api";
import { RepoInfo } from "../lib/types";

export default function Home() {
  const [repos, setRepos] = useState<RepoInfo[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchRepos = async () => {
      try {
        const response = await getRecentRepos();
        if (response.status === 200 && response.data && response.data.repos) {
          setRepos(response.data.repos);
        }
      } catch (error) {
        console.error("Failed to fetch recent repos:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchRepos();
  }, []);

  return (
    <div>
      <div className="pb-10 pt-10">
        <div className="mx-auto max-w-2xl px-4">
          <div className="flex flex-col gap-2 rounded-md border p-8">
            <h1 className="text-lg font-semibold">
              {">"} Print your favorite Git repositories as PDF
            </h1>
            <p className="py-2">
              Looking for a fun way to explore your favorite GitHub
              repositories? Tired of staring at the screen for hours on end? Or
              maybe want to print out a hard copy as a keepsake?
            </p>
            <p className="py-2">
              Simply sign in with your GitHub account and start printing{" "}
              <span className="text-white font-semibold">public</span>{" "}
              repositories in a beautiful, easy-to-read format.
            </p>
            <p className="py-2">
              It is currently in beta, so please be patient with us as we work.
              Feel free to request features or report bugs.
            </p>
            <p className="py-2 font-bold">
              Made by <a href="https://x.com/pliutau">@pliutau</a> with ❤️ | v0.2.0
            </p>
          </div>
        </div>
      </div>
      <div className="pb-10">{Form()}</div>
      {!loading && repos.length > 0 && (
        <div className="pb-10">
          <div className="mx-auto max-w-2xl px-4">
            <div className="flex flex-col text-center pb-4">
              <h2 className="text-lg font-semibold">Repositories of the day</h2>
              <span className="text-sm text-gray-400">
                Click to download the PDF
              </span>
            </div>
            <div className="mb-4 grid grid-cols-2 gap-2">
              {repos.map((repo) => (
                <a
                  key={repo.export_id}
                  className="cursor-pointer rounded-md border p-4"
                  href={`${process.env.NEXT_PUBLIC_API_ADDR}/files?export_id=${repo.export_id}&ext=pdf`}
                >
                  <div className="text-sm font-semibold">{repo.name}</div>
                  <div className="text-sm text-white">{repo.version}</div>
                  <div className="text-sm text-zinc-600">{repo.size}</div>
                </a>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
