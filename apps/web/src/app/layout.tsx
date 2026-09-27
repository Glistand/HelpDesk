import type { Metadata } from "next";
import { Inter } from "next/font/google";
import { ThemeProvider } from "@/components/theme-provider";
import { branding } from "@/lib/branding";
import "./globals.css";

const inter = Inter({
  variable: "--font-inter",
  subsets: ["latin", "cyrillic"],
});

export const dynamic = "force-dynamic";

export function generateMetadata(): Metadata {
  const { name } = branding();
  return {
    title: name,
    description: `Рабочее место операторов ${name}`,
  };
}

function themeInitScript(defaultTheme: "dark" | "light") {
  return `
(function () {
  try {
    var t = localStorage.getItem("helpdesk-theme");
    document.documentElement.setAttribute(
      "data-theme",
      t === "light" || t === "dark" ? t : "${defaultTheme}"
    );
  } catch (e) {
    document.documentElement.setAttribute("data-theme", "${defaultTheme}");
  }
})();
`;
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const { defaultTheme, style } = branding();

  return (
    <html
      lang="ru"
      className={`${inter.variable} h-full`}
      style={style}
      suppressHydrationWarning
    >
      <head>
        <script
          dangerouslySetInnerHTML={{ __html: themeInitScript(defaultTheme) }}
        />
      </head>
      <body className="min-h-full antialiased">
        <ThemeProvider defaultTheme={defaultTheme}>{children}</ThemeProvider>
      </body>
    </html>
  );
}
