import type { GlobalThemeOverrides } from 'naive-ui';

// 科幻主题色彩定义
export const sciFiColors = {
  dark: {
    // 主色调 - 赛博朋克青色
    primary: '#00f0ff',
    primaryHover: '#33f3ff',
    primaryPressed: '#00c8d4',
    primarySuppl: '#00a8b8',
    
    // 辅助色
    secondary: '#ff00ff',
    secondaryHover: '#ff33ff',
    accent: '#7c3aed',
    accentHover: '#8b5cf6',
    
    // 警告/错误/成功/信息
    success: '#00ff9f',
    warning: '#ffb800',
    error: '#ff3d71',
    info: '#00b8ff',
    
    // 背景色
    bgPrimary: '#0a0a0f',
    bgSecondary: '#12121a',
    bgTertiary: '#1a1a26',
    bgCard: '#16161f',
    bgCardHover: '#1e1e2a',
    bgModal: '#141420',
    bgOverlay: 'rgba(0, 0, 0, 0.8)',
    
    // 边框色
    border: '#2a2a3a',
    borderLight: '#3a3a4a',
    borderGlow: 'rgba(0, 240, 255, 0.3)',
    
    // 文字色
    textPrimary: '#ffffff',
    textSecondary: '#a0a0b0',
    textTertiary: '#6a6a7a',
    textDisabled: '#4a4a5a',
    
    // 特效色
    glowPrimary: '0 0 20px rgba(0, 240, 255, 0.5)',
    glowSecondary: '0 0 20px rgba(255, 0, 255, 0.3)',
    glowSuccess: '0 0 15px rgba(0, 255, 159, 0.4)',
    
    // 渐变
    gradientPrimary: 'linear-gradient(135deg, #00f0ff 0%, #7c3aed 100%)',
    gradientBg: 'linear-gradient(180deg, #0a0a0f 0%, #12121a 100%)',
    gradientCard: 'linear-gradient(145deg, #1a1a26 0%, #12121a 100%)'
  },
  light: {
    // 主色调 - 科技蓝
    primary: '#0066ff',
    primaryHover: '#3385ff',
    primaryPressed: '#0052cc',
    primarySuppl: '#004db8',
    
    // 辅助色
    secondary: '#9333ea',
    secondaryHover: '#a855f7',
    accent: '#6366f1',
    accentHover: '#818cf8',
    
    // 警告/错误/成功/信息
    success: '#10b981',
    warning: '#f59e0b',
    error: '#ef4444',
    info: '#3b82f6',
    
    // 背景色
    bgPrimary: '#f8faff',
    bgSecondary: '#ffffff',
    bgTertiary: '#f0f4ff',
    bgCard: '#ffffff',
    bgCardHover: '#f5f8ff',
    bgModal: '#ffffff',
    bgOverlay: 'rgba(0, 0, 0, 0.4)',
    
    // 边框色
    border: '#e0e6f0',
    borderLight: '#d0d8e8',
    borderGlow: 'rgba(0, 102, 255, 0.2)',
    
    // 文字色
    textPrimary: '#1a1a2e',
    textSecondary: '#4a5568',
    textTertiary: '#718096',
    textDisabled: '#a0aec0',
    
    // 特效色
    glowPrimary: '0 0 20px rgba(0, 102, 255, 0.25)',
    glowSecondary: '0 0 20px rgba(147, 51, 234, 0.2)',
    glowSuccess: '0 0 15px rgba(16, 185, 129, 0.25)',
    
    // 渐变
    gradientPrimary: 'linear-gradient(135deg, #0066ff 0%, #6366f1 100%)',
    gradientBg: 'linear-gradient(180deg, #f8faff 0%, #ffffff 100%)',
    gradientCard: 'linear-gradient(145deg, #ffffff 0%, #f5f8ff 100%)'
  }
};

// 暗黑主题覆盖
export const darkThemeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: sciFiColors.dark.primary,
    primaryColorHover: sciFiColors.dark.primaryHover,
    primaryColorPressed: sciFiColors.dark.primaryPressed,
    primaryColorSuppl: sciFiColors.dark.primarySuppl,
    
    infoColor: sciFiColors.dark.info,
    infoColorHover: '#33c7ff',
    infoColorPressed: '#0090cc',
    
    successColor: sciFiColors.dark.success,
    successColorHover: '#33ffa8',
    successColorPressed: '#00cc7f',
    
    warningColor: sciFiColors.dark.warning,
    warningColorHover: '#ffc633',
    warningColorPressed: '#cc9300',
    
    errorColor: sciFiColors.dark.error,
    errorColorHover: '#ff637d',
    errorColorPressed: '#cc315b',
    
    textColorBase: sciFiColors.dark.textPrimary,
    textColor1: sciFiColors.dark.textPrimary,
    textColor2: sciFiColors.dark.textSecondary,
    textColor3: sciFiColors.dark.textTertiary,
    textColorDisabled: sciFiColors.dark.textDisabled,
    
    placeholderColor: sciFiColors.dark.textTertiary,
    placeholderColorDisabled: sciFiColors.dark.textDisabled,
    
    iconColor: sciFiColors.dark.textSecondary,
    iconColorHover: sciFiColors.dark.primary,
    iconColorPressed: sciFiColors.dark.primaryPressed,
    iconColorDisabled: sciFiColors.dark.textDisabled,
    
    dividerColor: sciFiColors.dark.border,
    borderColor: sciFiColors.dark.border,
    
    closeIconColor: sciFiColors.dark.textSecondary,
    closeIconColorHover: sciFiColors.dark.textPrimary,
    closeIconColorPressed: sciFiColors.dark.primary,
    
    clearColor: sciFiColors.dark.textTertiary,
    clearColorHover: sciFiColors.dark.textSecondary,
    clearColorPressed: sciFiColors.dark.textPrimary,
    
    scrollbarColor: 'rgba(0, 240, 255, 0.2)',
    scrollbarColorHover: 'rgba(0, 240, 255, 0.4)',
    
    progressRailColor: sciFiColors.dark.border,
    
    bodyColor: sciFiColors.dark.bgPrimary,
    cardColor: sciFiColors.dark.bgCard,
    modalColor: sciFiColors.dark.bgModal,
    popoverColor: sciFiColors.dark.bgCard,
    tableColor: sciFiColors.dark.bgCard,
    inputColor: sciFiColors.dark.bgTertiary,
    inputColorDisabled: sciFiColors.dark.bgSecondary,
    actionColor: sciFiColors.dark.bgTertiary,
    tagColor: sciFiColors.dark.bgTertiary,
    avatarColor: sciFiColors.dark.bgTertiary,
    invertedColor: sciFiColors.dark.bgCard,
    
    hoverColor: 'rgba(0, 240, 255, 0.08)',
    pressedColor: 'rgba(0, 240, 255, 0.12)',
    
    boxShadow1: '0 4px 20px rgba(0, 0, 0, 0.4)',
    boxShadow2: '0 8px 30px rgba(0, 0, 0, 0.5)',
    boxShadow3: '0 12px 40px rgba(0, 0, 0, 0.6)',
    
    borderRadius: '8px',
    borderRadiusSmall: '4px',
    
    fontFamily: '"Inter", "SF Pro Display", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
    fontFamilyMono: '"JetBrains Mono", "Fira Code", "SF Mono", Monaco, Consolas, monospace'
  },
  Button: {
    textColorPrimary: '#000000',
    textColorHoverPrimary: '#000000',
    textColorPressedPrimary: '#000000',
    textColorFocusPrimary: '#000000',
    colorPrimary: sciFiColors.dark.primary,
    colorHoverPrimary: sciFiColors.dark.primaryHover,
    colorPressedPrimary: sciFiColors.dark.primaryPressed,
    colorFocusPrimary: sciFiColors.dark.primary,
    borderPrimary: '1px solid transparent',
    borderHoverPrimary: '1px solid transparent',
    borderPressedPrimary: '1px solid transparent',
    borderFocusPrimary: '1px solid transparent',
    rippleColorPrimary: sciFiColors.dark.primary,
    
    textColorGhost: sciFiColors.dark.primary,
    textColorGhostHover: sciFiColors.dark.primaryHover,
    textColorGhostPressed: sciFiColors.dark.primaryPressed,
    borderGhost: `1px solid ${sciFiColors.dark.primary}`,
    borderGhostHover: `1px solid ${sciFiColors.dark.primaryHover}`,
    borderGhostPressed: `1px solid ${sciFiColors.dark.primaryPressed}`,
    
    borderRadiusMedium: '8px',
    borderRadiusSmall: '6px',
    borderRadiusLarge: '10px',
    
    fontWeightStrong: '600'
  },
  Card: {
    color: sciFiColors.dark.bgCard,
    colorModal: sciFiColors.dark.bgModal,
    colorPopover: sciFiColors.dark.bgCard,
    colorTarget: sciFiColors.dark.bgCard,
    colorEmbedded: sciFiColors.dark.bgSecondary,
    borderColor: sciFiColors.dark.border,
    borderRadius: '12px',
    boxShadow: '0 4px 24px rgba(0, 0, 0, 0.3)',
    titleFontSizeMedium: '18px',
    titleFontWeight: '600'
  },
  Input: {
    color: sciFiColors.dark.bgTertiary,
    colorFocus: sciFiColors.dark.bgTertiary,
    colorDisabled: sciFiColors.dark.bgSecondary,
    border: `1px solid ${sciFiColors.dark.border}`,
    borderHover: `1px solid ${sciFiColors.dark.primary}`,
    borderFocus: `1px solid ${sciFiColors.dark.primary}`,
    borderDisabled: `1px solid ${sciFiColors.dark.border}`,
    boxShadowFocus: `0 0 0 2px rgba(0, 240, 255, 0.2)`,
    caretColor: sciFiColors.dark.primary,
    borderRadius: '8px',
    heightMedium: '40px'
  },
  Select: {
    peers: {
      InternalSelection: {
        color: sciFiColors.dark.bgTertiary,
        colorActive: sciFiColors.dark.bgTertiary,
        colorDisabled: sciFiColors.dark.bgSecondary,
        border: `1px solid ${sciFiColors.dark.border}`,
        borderHover: `1px solid ${sciFiColors.dark.primary}`,
        borderActive: `1px solid ${sciFiColors.dark.primary}`,
        borderFocus: `1px solid ${sciFiColors.dark.primary}`,
        boxShadowFocus: `0 0 0 2px rgba(0, 240, 255, 0.2)`,
        boxShadowActive: `0 0 0 2px rgba(0, 240, 255, 0.2)`,
        borderRadius: '8px',
        textColor: sciFiColors.dark.textPrimary,
        placeholderColor: sciFiColors.dark.textTertiary,
        arrowColor: sciFiColors.dark.textSecondary
      },
      InternalSelectMenu: {
        color: sciFiColors.dark.bgCard,
        optionColorPending: 'rgba(0, 240, 255, 0.1)',
        optionColorActive: 'rgba(0, 240, 255, 0.15)',
        optionTextColorActive: sciFiColors.dark.primary,
        optionCheckColor: sciFiColors.dark.primary,
        borderRadius: '8px'
      }
    }
  },
  Menu: {
    color: 'transparent',
    itemColorHover: 'rgba(0, 240, 255, 0.08)',
    itemColorActive: 'rgba(0, 240, 255, 0.12)',
    itemColorActiveHover: 'rgba(0, 240, 255, 0.15)',
    itemColorActiveCollapsed: 'rgba(0, 240, 255, 0.12)',
    itemTextColor: sciFiColors.dark.textSecondary,
    itemTextColorHover: sciFiColors.dark.textPrimary,
    itemTextColorActive: sciFiColors.dark.primary,
    itemTextColorActiveHover: sciFiColors.dark.primary,
    itemTextColorChildActive: sciFiColors.dark.primary,
    itemIconColor: sciFiColors.dark.textSecondary,
    itemIconColorHover: sciFiColors.dark.textPrimary,
    itemIconColorActive: sciFiColors.dark.primary,
    itemIconColorActiveHover: sciFiColors.dark.primary,
    itemIconColorCollapsed: sciFiColors.dark.textSecondary,
    arrowColor: sciFiColors.dark.textSecondary,
    arrowColorHover: sciFiColors.dark.textPrimary,
    arrowColorActive: sciFiColors.dark.primary,
    arrowColorChildActive: sciFiColors.dark.primary,
    borderRadius: '8px',
    itemHeight: '44px'
  },
  Modal: {
    color: sciFiColors.dark.bgModal,
    textColor: sciFiColors.dark.textPrimary,
    borderRadius: '16px',
    boxShadow: '0 20px 60px rgba(0, 0, 0, 0.5)'
  },
  Message: {
    borderRadius: '10px',
    boxShadow: '0 8px 24px rgba(0, 0, 0, 0.4)'
  },
  Tag: {
    color: sciFiColors.dark.bgTertiary,
    textColor: sciFiColors.dark.textSecondary,
    border: `1px solid ${sciFiColors.dark.border}`,
    borderRadius: '6px',
    colorPrimary: 'rgba(0, 240, 255, 0.15)',
    textColorPrimary: sciFiColors.dark.primary,
    borderPrimary: `1px solid ${sciFiColors.dark.primary}`,
    closeIconColor: sciFiColors.dark.textTertiary,
    closeIconColorHover: sciFiColors.dark.textPrimary
  },
  Avatar: {
    color: sciFiColors.dark.bgTertiary,
    borderRadius: '50%'
  },
  Tabs: {
    tabTextColorLine: sciFiColors.dark.textSecondary,
    tabTextColorActiveLine: sciFiColors.dark.primary,
    tabTextColorHoverLine: sciFiColors.dark.textPrimary,
    barColor: sciFiColors.dark.primary,
    tabBorderColor: sciFiColors.dark.border,
    tabGapMediumLine: '24px',
    tabPaddingMediumLine: '12px 0',
    tabFontWeightActive: '600'
  },
  Scrollbar: {
    color: 'rgba(0, 240, 255, 0.15)',
    colorHover: 'rgba(0, 240, 255, 0.3)',
    width: '8px',
    borderRadius: '4px'
  },
  Drawer: {
    color: sciFiColors.dark.bgSecondary,
    textColor: sciFiColors.dark.textPrimary,
    borderRadius: '16px 0 0 16px',
    boxShadow: '-20px 0 60px rgba(0, 0, 0, 0.4)'
  },
  Switch: {
    railColor: sciFiColors.dark.border,
    railColorActive: sciFiColors.dark.primary,
    buttonColor: sciFiColors.dark.textPrimary,
    boxShadowFocus: `0 0 0 2px rgba(0, 240, 255, 0.2)`
  },
  Slider: {
    fillColor: sciFiColors.dark.primary,
    fillColorHover: sciFiColors.dark.primaryHover,
    railColor: sciFiColors.dark.border,
    handleColor: sciFiColors.dark.primary,
    dotBorderActive: `2px solid ${sciFiColors.dark.primary}`
  },
  Progress: {
    fillColor: sciFiColors.dark.primary,
    railColor: sciFiColors.dark.border,
    textColorCircle: sciFiColors.dark.textPrimary
  },
  List: {
    color: 'transparent',
    colorHover: 'rgba(0, 240, 255, 0.05)',
    borderColor: sciFiColors.dark.border,
    textColor: sciFiColors.dark.textPrimary
  },
  Empty: {
    textColor: sciFiColors.dark.textTertiary,
    iconColor: sciFiColors.dark.textTertiary
  },
  Tooltip: {
    color: sciFiColors.dark.bgCard,
    textColor: sciFiColors.dark.textPrimary,
    borderRadius: '8px',
    boxShadow: '0 8px 24px rgba(0, 0, 0, 0.4)'
  },
  Popover: {
    color: sciFiColors.dark.bgCard,
    textColor: sciFiColors.dark.textPrimary,
    borderRadius: '12px',
    boxShadow: '0 12px 32px rgba(0, 0, 0, 0.4)'
  },
  Dropdown: {
    color: sciFiColors.dark.bgCard,
    optionColorHover: 'rgba(0, 240, 255, 0.1)',
    optionColorActive: 'rgba(0, 240, 255, 0.15)',
    optionTextColor: sciFiColors.dark.textSecondary,
    optionTextColorHover: sciFiColors.dark.textPrimary,
    optionTextColorActive: sciFiColors.dark.primary,
    borderRadius: '10px',
    boxShadow: '0 12px 32px rgba(0, 0, 0, 0.4)'
  }
};

// 明亮主题覆盖
export const lightThemeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: sciFiColors.light.primary,
    primaryColorHover: sciFiColors.light.primaryHover,
    primaryColorPressed: sciFiColors.light.primaryPressed,
    primaryColorSuppl: sciFiColors.light.primarySuppl,
    
    infoColor: sciFiColors.light.info,
    infoColorHover: '#5a9efc',
    infoColorPressed: '#2563eb',
    
    successColor: sciFiColors.light.success,
    successColorHover: '#34d399',
    successColorPressed: '#059669',
    
    warningColor: sciFiColors.light.warning,
    warningColorHover: '#fbbf24',
    warningColorPressed: '#d97706',
    
    errorColor: sciFiColors.light.error,
    errorColorHover: '#f87171',
    errorColorPressed: '#dc2626',
    
    textColorBase: sciFiColors.light.textPrimary,
    textColor1: sciFiColors.light.textPrimary,
    textColor2: sciFiColors.light.textSecondary,
    textColor3: sciFiColors.light.textTertiary,
    textColorDisabled: sciFiColors.light.textDisabled,
    
    placeholderColor: sciFiColors.light.textTertiary,
    placeholderColorDisabled: sciFiColors.light.textDisabled,
    
    iconColor: sciFiColors.light.textSecondary,
    iconColorHover: sciFiColors.light.primary,
    iconColorPressed: sciFiColors.light.primaryPressed,
    iconColorDisabled: sciFiColors.light.textDisabled,
    
    dividerColor: sciFiColors.light.border,
    borderColor: sciFiColors.light.border,
    
    closeIconColor: sciFiColors.light.textSecondary,
    closeIconColorHover: sciFiColors.light.textPrimary,
    closeIconColorPressed: sciFiColors.light.primary,
    
    clearColor: sciFiColors.light.textTertiary,
    clearColorHover: sciFiColors.light.textSecondary,
    clearColorPressed: sciFiColors.light.textPrimary,
    
    scrollbarColor: 'rgba(0, 102, 255, 0.15)',
    scrollbarColorHover: 'rgba(0, 102, 255, 0.3)',
    
    progressRailColor: sciFiColors.light.border,
    
    bodyColor: sciFiColors.light.bgPrimary,
    cardColor: sciFiColors.light.bgCard,
    modalColor: sciFiColors.light.bgModal,
    popoverColor: sciFiColors.light.bgCard,
    tableColor: sciFiColors.light.bgCard,
    inputColor: sciFiColors.light.bgTertiary,
    inputColorDisabled: sciFiColors.light.bgSecondary,
    actionColor: sciFiColors.light.bgTertiary,
    tagColor: sciFiColors.light.bgTertiary,
    avatarColor: sciFiColors.light.bgTertiary,
    invertedColor: sciFiColors.light.bgCard,
    
    hoverColor: 'rgba(0, 102, 255, 0.06)',
    pressedColor: 'rgba(0, 102, 255, 0.1)',
    
    boxShadow1: '0 4px 20px rgba(0, 0, 0, 0.08)',
    boxShadow2: '0 8px 30px rgba(0, 0, 0, 0.1)',
    boxShadow3: '0 12px 40px rgba(0, 0, 0, 0.12)',
    
    borderRadius: '8px',
    borderRadiusSmall: '4px',
    
    fontFamily: '"Inter", "SF Pro Display", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
    fontFamilyMono: '"JetBrains Mono", "Fira Code", "SF Mono", Monaco, Consolas, monospace'
  },
  Button: {
    textColorPrimary: '#ffffff',
    textColorHoverPrimary: '#ffffff',
    textColorPressedPrimary: '#ffffff',
    textColorFocusPrimary: '#ffffff',
    colorPrimary: sciFiColors.light.primary,
    colorHoverPrimary: sciFiColors.light.primaryHover,
    colorPressedPrimary: sciFiColors.light.primaryPressed,
    colorFocusPrimary: sciFiColors.light.primary,
    borderPrimary: '1px solid transparent',
    borderHoverPrimary: '1px solid transparent',
    borderPressedPrimary: '1px solid transparent',
    borderFocusPrimary: '1px solid transparent',
    rippleColorPrimary: sciFiColors.light.primary,
    
    textColorGhost: sciFiColors.light.primary,
    textColorGhostHover: sciFiColors.light.primaryHover,
    textColorGhostPressed: sciFiColors.light.primaryPressed,
    borderGhost: `1px solid ${sciFiColors.light.primary}`,
    borderGhostHover: `1px solid ${sciFiColors.light.primaryHover}`,
    borderGhostPressed: `1px solid ${sciFiColors.light.primaryPressed}`,
    
    borderRadiusMedium: '8px',
    borderRadiusSmall: '6px',
    borderRadiusLarge: '10px',
    
    fontWeightStrong: '600'
  },
  Card: {
    color: sciFiColors.light.bgCard,
    colorModal: sciFiColors.light.bgModal,
    colorPopover: sciFiColors.light.bgCard,
    colorTarget: sciFiColors.light.bgCard,
    colorEmbedded: sciFiColors.light.bgSecondary,
    borderColor: sciFiColors.light.border,
    borderRadius: '12px',
    boxShadow: '0 4px 24px rgba(0, 0, 0, 0.08)',
    titleFontSizeMedium: '18px',
    titleFontWeight: '600'
  },
  Input: {
    color: sciFiColors.light.bgSecondary,
    colorFocus: sciFiColors.light.bgSecondary,
    colorDisabled: sciFiColors.light.bgTertiary,
    border: `1px solid ${sciFiColors.light.border}`,
    borderHover: `1px solid ${sciFiColors.light.primary}`,
    borderFocus: `1px solid ${sciFiColors.light.primary}`,
    borderDisabled: `1px solid ${sciFiColors.light.border}`,
    boxShadowFocus: `0 0 0 2px rgba(0, 102, 255, 0.15)`,
    caretColor: sciFiColors.light.primary,
    borderRadius: '8px',
    heightMedium: '40px'
  },
  Select: {
    peers: {
      InternalSelection: {
        color: sciFiColors.light.bgSecondary,
        colorActive: sciFiColors.light.bgSecondary,
        colorDisabled: sciFiColors.light.bgTertiary,
        border: `1px solid ${sciFiColors.light.border}`,
        borderHover: `1px solid ${sciFiColors.light.primary}`,
        borderActive: `1px solid ${sciFiColors.light.primary}`,
        borderFocus: `1px solid ${sciFiColors.light.primary}`,
        boxShadowFocus: `0 0 0 2px rgba(0, 102, 255, 0.15)`,
        boxShadowActive: `0 0 0 2px rgba(0, 102, 255, 0.15)`,
        borderRadius: '8px',
        textColor: sciFiColors.light.textPrimary,
        placeholderColor: sciFiColors.light.textTertiary,
        arrowColor: sciFiColors.light.textSecondary
      },
      InternalSelectMenu: {
        color: sciFiColors.light.bgCard,
        optionColorPending: 'rgba(0, 102, 255, 0.08)',
        optionColorActive: 'rgba(0, 102, 255, 0.12)',
        optionTextColorActive: sciFiColors.light.primary,
        optionCheckColor: sciFiColors.light.primary,
        borderRadius: '8px'
      }
    }
  },
  Menu: {
    color: 'transparent',
    itemColorHover: 'rgba(0, 102, 255, 0.06)',
    itemColorActive: 'rgba(0, 102, 255, 0.1)',
    itemColorActiveHover: 'rgba(0, 102, 255, 0.12)',
    itemColorActiveCollapsed: 'rgba(0, 102, 255, 0.1)',
    itemTextColor: sciFiColors.light.textSecondary,
    itemTextColorHover: sciFiColors.light.textPrimary,
    itemTextColorActive: sciFiColors.light.primary,
    itemTextColorActiveHover: sciFiColors.light.primary,
    itemTextColorChildActive: sciFiColors.light.primary,
    itemIconColor: sciFiColors.light.textSecondary,
    itemIconColorHover: sciFiColors.light.textPrimary,
    itemIconColorActive: sciFiColors.light.primary,
    itemIconColorActiveHover: sciFiColors.light.primary,
    itemIconColorCollapsed: sciFiColors.light.textSecondary,
    arrowColor: sciFiColors.light.textSecondary,
    arrowColorHover: sciFiColors.light.textPrimary,
    arrowColorActive: sciFiColors.light.primary,
    arrowColorChildActive: sciFiColors.light.primary,
    borderRadius: '8px',
    itemHeight: '44px'
  },
  Modal: {
    color: sciFiColors.light.bgModal,
    textColor: sciFiColors.light.textPrimary,
    borderRadius: '16px',
    boxShadow: '0 20px 60px rgba(0, 0, 0, 0.15)'
  },
  Message: {
    borderRadius: '10px',
    boxShadow: '0 8px 24px rgba(0, 0, 0, 0.1)'
  },
  Tag: {
    color: sciFiColors.light.bgTertiary,
    textColor: sciFiColors.light.textSecondary,
    border: `1px solid ${sciFiColors.light.border}`,
    borderRadius: '6px',
    colorPrimary: 'rgba(0, 102, 255, 0.1)',
    textColorPrimary: sciFiColors.light.primary,
    borderPrimary: `1px solid ${sciFiColors.light.primary}`,
    closeIconColor: sciFiColors.light.textTertiary,
    closeIconColorHover: sciFiColors.light.textPrimary
  },
  Avatar: {
    color: sciFiColors.light.bgTertiary,
    borderRadius: '50%'
  },
  Tabs: {
    tabTextColorLine: sciFiColors.light.textSecondary,
    tabTextColorActiveLine: sciFiColors.light.primary,
    tabTextColorHoverLine: sciFiColors.light.textPrimary,
    barColor: sciFiColors.light.primary,
    tabBorderColor: sciFiColors.light.border,
    tabGapMediumLine: '24px',
    tabPaddingMediumLine: '12px 0',
    tabFontWeightActive: '600'
  },
  Scrollbar: {
    color: 'rgba(0, 102, 255, 0.1)',
    colorHover: 'rgba(0, 102, 255, 0.2)',
    width: '8px',
    borderRadius: '4px'
  },
  Drawer: {
    color: sciFiColors.light.bgSecondary,
    textColor: sciFiColors.light.textPrimary,
    borderRadius: '16px 0 0 16px',
    boxShadow: '-20px 0 60px rgba(0, 0, 0, 0.1)'
  },
  Switch: {
    railColor: sciFiColors.light.border,
    railColorActive: sciFiColors.light.primary,
    buttonColor: '#ffffff',
    boxShadowFocus: `0 0 0 2px rgba(0, 102, 255, 0.15)`
  },
  Slider: {
    fillColor: sciFiColors.light.primary,
    fillColorHover: sciFiColors.light.primaryHover,
    railColor: sciFiColors.light.border,
    handleColor: sciFiColors.light.primary,
    dotBorderActive: `2px solid ${sciFiColors.light.primary}`
  },
  Progress: {
    fillColor: sciFiColors.light.primary,
    railColor: sciFiColors.light.border,
    textColorCircle: sciFiColors.light.textPrimary
  },
  List: {
    color: 'transparent',
    colorHover: 'rgba(0, 102, 255, 0.04)',
    borderColor: sciFiColors.light.border,
    textColor: sciFiColors.light.textPrimary
  },
  Empty: {
    textColor: sciFiColors.light.textTertiary,
    iconColor: sciFiColors.light.textTertiary
  },
  Tooltip: {
    color: sciFiColors.light.bgCard,
    textColor: sciFiColors.light.textPrimary,
    borderRadius: '8px',
    boxShadow: '0 8px 24px rgba(0, 0, 0, 0.1)'
  },
  Popover: {
    color: sciFiColors.light.bgCard,
    textColor: sciFiColors.light.textPrimary,
    borderRadius: '12px',
    boxShadow: '0 12px 32px rgba(0, 0, 0, 0.1)'
  },
  Dropdown: {
    color: sciFiColors.light.bgCard,
    optionColorHover: 'rgba(0, 102, 255, 0.06)',
    optionColorActive: 'rgba(0, 102, 255, 0.1)',
    optionTextColor: sciFiColors.light.textSecondary,
    optionTextColorHover: sciFiColors.light.textPrimary,
    optionTextColorActive: sciFiColors.light.primary,
    borderRadius: '10px',
    boxShadow: '0 12px 32px rgba(0, 0, 0, 0.1)'
  }
};
