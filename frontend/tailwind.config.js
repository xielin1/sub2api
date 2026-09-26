/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // 主色调 - 与首页一致的陶土暖橙，500 为首页主按钮色 #d66b4d
        primary: {
          50: '#fcf5f1',
          100: '#f9e7df',
          200: '#f2cdbd',
          300: '#e9ad96',
          400: '#e0886b',
          500: '#d66b4d',
          600: '#c25f43',
          700: '#a24a33',
          800: '#843d2c',
          900: '#6c3427',
          950: '#3a1912'
        },
        // 浅色模式中性色 - 暖米灰，替换默认偏蓝的冷灰
        gray: {
          50: '#f7f5ef',
          100: '#f0ece3',
          200: '#e4dfd3',
          300: '#d2ccbe',
          400: '#a8a193',
          500: '#7d776a',
          600: '#5c574c',
          700: '#45413a',
          800: '#2c2a24',
          900: '#1f1e1a',
          950: '#121310'
        },
        slate: {
          50: '#f7f5ef',
          100: '#f0ece3',
          200: '#e4dfd3',
          300: '#d2ccbe',
          400: '#a8a193',
          500: '#7d776a',
          600: '#5c574c',
          700: '#45413a',
          800: '#2c2a24',
          900: '#1f1e1a',
          950: '#121310'
        },
        // 辅助色 - 暖灰
        accent: {
          50: '#f7f5ef',
          100: '#f0ece3',
          200: '#e4dfd3',
          300: '#d2ccbe',
          400: '#a8a193',
          500: '#7d776a',
          600: '#5c574c',
          700: '#45413a',
          800: '#2c2a24',
          900: '#1f1e1a',
          950: '#121310'
        },
        // 深色模式背景 - 与首页深色一致的暖黑（900/950 即首页 #1b1c19 / #121310）
        dark: {
          50: '#f7f5ef',
          100: '#f3f1ea',
          200: '#e5e1d6',
          300: '#c9c4b6',
          400: '#a39d8f',
          500: '#7d776a',
          600: '#4d4a40',
          700: '#3a382f',
          800: '#26261f',
          900: '#1b1c19',
          950: '#121310'
        }
      },
      fontFamily: {
        sans: [
          'system-ui',
          '-apple-system',
          'BlinkMacSystemFont',
          'Segoe UI',
          'Roboto',
          'Helvetica Neue',
          'Arial',
          'PingFang SC',
          'Hiragino Sans GB',
          'Microsoft YaHei',
          'sans-serif'
        ],
        mono: ['ui-monospace', 'SFMono-Regular', 'Menlo', 'Monaco', 'Consolas', 'monospace']
      },
      boxShadow: {
        glass: '0 8px 32px rgba(0, 0, 0, 0.08)',
        'glass-sm': '0 4px 16px rgba(0, 0, 0, 0.06)',
        glow: '0 0 20px rgba(214, 107, 77, 0.25)',
        'glow-lg': '0 0 40px rgba(214, 107, 77, 0.35)',
        card: '0 1px 3px rgba(0, 0, 0, 0.04), 0 1px 2px rgba(0, 0, 0, 0.06)',
        'card-hover': '0 10px 40px rgba(0, 0, 0, 0.08)',
        'inner-glow': 'inset 0 1px 0 rgba(255, 255, 255, 0.1)'
      },
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'gradient-primary': 'linear-gradient(135deg, #d66b4d 0%, #c25f43 100%)',
        'gradient-dark': 'linear-gradient(135deg, #26261f 0%, #1b1c19 100%)',
        'gradient-glass':
          'linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0.05) 100%)',
        'mesh-gradient':
          'radial-gradient(at 40% 20%, rgba(214, 107, 77, 0.10) 0px, transparent 50%), radial-gradient(at 80% 0%, rgba(226, 176, 138, 0.10) 0px, transparent 50%), radial-gradient(at 0% 50%, rgba(214, 107, 77, 0.06) 0px, transparent 50%)'
      },
      animation: {
        'fade-in': 'fadeIn 0.3s ease-out',
        'slide-up': 'slideUp 0.3s ease-out',
        'slide-down': 'slideDown 0.3s ease-out',
        'slide-in-right': 'slideInRight 0.3s ease-out',
        'scale-in': 'scaleIn 0.2s ease-out',
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        shimmer: 'shimmer 2s linear infinite',
        glow: 'glow 2s ease-in-out infinite alternate'
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' }
        },
        slideUp: {
          '0%': { opacity: '0', transform: 'translateY(10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' }
        },
        slideDown: {
          '0%': { opacity: '0', transform: 'translateY(-10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' }
        },
        slideInRight: {
          '0%': { opacity: '0', transform: 'translateX(20px)' },
          '100%': { opacity: '1', transform: 'translateX(0)' }
        },
        scaleIn: {
          '0%': { opacity: '0', transform: 'scale(0.95)' },
          '100%': { opacity: '1', transform: 'scale(1)' }
        },
        shimmer: {
          '0%': { backgroundPosition: '-200% 0' },
          '100%': { backgroundPosition: '200% 0' }
        },
        glow: {
          '0%': { boxShadow: '0 0 20px rgba(214, 107, 77, 0.25)' },
          '100%': { boxShadow: '0 0 30px rgba(214, 107, 77, 0.4)' }
        }
      },
      backdropBlur: {
        xs: '2px'
      },
      borderRadius: {
        '4xl': '2rem'
      }
    }
  },
  plugins: []
}
