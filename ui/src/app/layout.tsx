import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "gitprint.me",
  description: "Print your favorite Git repositories as PDF.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark">
      <body>
        <header className="fixed top-0 z-50 flex h-16 w-full items-center justify-center border-b">
          <a rel="nofollow" href="/" className="font-bold p-8">
            gitprint.me - Print your favorite Git repositories as PDF.
          </a>
        </header>
        <main className="flex h-dvh flex-1 flex-col pt-16">{children}</main>
      </body>
    </html>
  );
}
