import React from 'react';
import { Bot } from 'lucide-react';

interface AuthLayoutProps {
  children: React.ReactNode;
  title: string;
  subtitle: string;
}

export function AuthLayout({ children, title, subtitle }: AuthLayoutProps) {
  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-950 text-slate-100 font-sans antialiased relative overflow-hidden px-4">
      {/* Premium Background Mesh and Radial Glows */}
      <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-indigo-950/30 via-slate-950 to-slate-950 -z-10" />
      <div className="absolute top-1/4 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[500px] h-[500px] bg-indigo-500/10 rounded-full blur-[120px] -z-10 animate-pulse pointer-events-none" />
      
      {/* Decorative Grid Gridline */}
      <div className="absolute inset-0 bg-[linear-gradient(to_right,#0f172a_1px,transparent_1px),linear-gradient(to_bottom,#0f172a_1px,transparent_1px)] bg-[size:4rem_4rem] [mask-image:radial-gradient(ellipse_60%_50%_at_50%_50%,#000_70%,transparent_100%)] opacity-35 -z-10 pointer-events-none" />

      <div className="w-full max-w-md">
        {/* Core Logo Branding */}
        <div className="flex flex-col items-center mb-8 text-center space-y-3">
          <div className="p-3 rounded-2xl bg-indigo-600/10 text-indigo-400 border border-indigo-500/20 shadow-lg shadow-indigo-500/5 animate-bounce">
            <Bot size={28} />
          </div>
          <div>
            <h1 className="text-xl font-bold tracking-tight text-white">{title}</h1>
            <p className="text-xs text-slate-400 mt-1">{subtitle}</p>
          </div>
        </div>

        {/* Glassmorphic Panel Card */}
        <div className="bg-slate-900/40 border border-slate-800/80 rounded-2xl p-6 md:p-8 backdrop-blur-md shadow-2xl relative">
          {/* Subtle top indicator border */}
          <div className="absolute top-0 left-10 right-10 h-[1px] bg-gradient-to-right from-transparent via-indigo-500 to-transparent opacity-60" />
          
          {children}
        </div>
      </div>
    </div>
  );
}
