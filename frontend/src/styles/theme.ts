import type { GlobalThemeOverrides } from 'naive-ui';

// Google AI Studio 风格主题色彩定义
export const museColors = {
  dark: {
    // 主色调 - 明亮电光蓝
    primary: '#4DA8FF',
    primaryHover: '#6FBCFF',
    primaryPressed: '#2B8CED',
    primarySuppl: '#3D9AFF',

    // 语义色
    success: '#5EE8A0',
    warning: '#FFD166',
    error: '#FF6B6B',
    info: '#4DA8FF',

    // 背景色 - 深邃太空感
    bgPrimary: '#0D1117',
    bgSecondary: '#161B22',
    bgTertiary: '#1C2333',
    bgCard: '#1C2333',
    bgCardHover: '#242D3D',
    bgModal: '#1C2333',
    bgOverlay: 'rgba(0, 0, 0, 0.7)',

    // 边框色 - 微光感
    border: 'rgba(77, 168, 255, 0.12)',
    borderLight: 'rgba(77, 168, 255, 0.2)',

    // 文字色
    textPrimary: '#E6EDF3',
    textSecondary: '#8B949E',
    textTertiary: '#6E7681',
    textDisabled: '#484F58',

    // hover/active
    bgHover: 'rgba(77, 168, 255, 0.08)',
    bgActive: 'rgba(77, 168, 255, 0.14)'
  },
  light: {
    // 主色调 - 鲜艳蓝
    primary: '#2563EB',
    primaryHover: '#1D4ED8',
    primaryPressed: '#1E40AF',
    primarySuppl: '#3B82F6',

    // 语义色
    success: '#16A34A',
    warning: '#F59E0B',
    error: '#DC2626',
    info: '#2563EB',

    // 背景色
    bgPrimary: '#FFFFFF',
    bgSecondary: '#F8FAFC',
    bgTertiary: '#F1F5F9',
    bgCard: '#FFFFFF',
    bgCardHover: '#F1F5F9',
    bgModal: '#FFFFFF',
    bgOverlay: 'rgba(0, 0, 0, 0.4)',

    // 边框色
    border: '#E2E8F0',
    borderLight: '#CBD5E1',

    // 文字色
    textPrimary: '#0F172A',
    textSecondary: '#475569',
    textTertiary: '#94A3B8',
    textDisabled: '#CBD5E1',

    // hover/active
    bgHover: '#F1F5F9',
    bgActive: '#EFF6FF'
  }
};

// 暗黑主题覆盖
export const darkThemeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: museColors.dark.primary,
    primaryColorHover: museColors.dark.primaryHover,
    primaryColorPressed: museColors.dark.primaryPressed,
    primaryColorSuppl: museColors.dark.primarySuppl,

    infoColor: museColors.dark.info,
    infoColorHover: '#6FBCFF',
    infoColorPressed: '#2B8CED',

    successColor: museColors.dark.success,
    successColorHover: '#7CEEB5',
    successColorPressed: '#3DD989',

    warningColor: museColors.dark.warning,
    warningColorHover: '#FFE08A',
    warningColorPressed: '#F5C230',

    errorColor: museColors.dark.error,
    errorColorHover: '#FF8A8A',
    errorColorPressed: '#E04D4D',

    textColorBase: museColors.dark.textPrimary,
    textColor1: museColors.dark.textPrimary,
    textColor2: museColors.dark.textSecondary,
    textColor3: museColors.dark.textTertiary,
    textColorDisabled: museColors.dark.textDisabled,

    placeholderColor: museColors.dark.textTertiary,
    placeholderColorDisabled: museColors.dark.textDisabled,

    iconColor: museColors.dark.textSecondary,
    iconColorHover: museColors.dark.textPrimary,
    iconColorPressed: museColors.dark.primary,
    iconColorDisabled: museColors.dark.textDisabled,

    dividerColor: museColors.dark.border,
    borderColor: museColors.dark.border,

    closeIconColor: museColors.dark.textSecondary,
    closeIconColorHover: museColors.dark.textPrimary,
    closeIconColorPressed: museColors.dark.primary,

    clearColor: museColors.dark.textTertiary,
    clearColorHover: museColors.dark.textSecondary,
    clearColorPressed: museColors.dark.textPrimary,

    scrollbarColor: 'rgba(77, 168, 255, 0.12)',
    scrollbarColorHover: 'rgba(77, 168, 255, 0.24)',

    progressRailColor: museColors.dark.border,

    bodyColor: museColors.dark.bgPrimary,
    cardColor: museColors.dark.bgCard,
    modalColor: museColors.dark.bgModal,
    popoverColor: museColors.dark.bgCard,
    tableColor: museColors.dark.bgCard,
    inputColor: museColors.dark.bgTertiary,
    inputColorDisabled: museColors.dark.bgSecondary,
    actionColor: museColors.dark.bgTertiary,
    tagColor: museColors.dark.bgTertiary,
    avatarColor: museColors.dark.bgTertiary,
    invertedColor: museColors.dark.bgCard,

    hoverColor: 'rgba(77, 168, 255, 0.06)',
    pressedColor: 'rgba(77, 168, 255, 0.12)',

    boxShadow1: '0 1px 3px rgba(0, 0, 0, 0.4)',
    boxShadow2: '0 4px 12px rgba(0, 0, 0, 0.5)',
    boxShadow3: '0 8px 24px rgba(0, 0, 0, 0.6)',

    borderRadius: '8px',
    borderRadiusSmall: '4px',

    fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Noto Sans SC", sans-serif',
    fontFamilyMono: '"JetBrains Mono", "Fira Code", "Source Code Pro", monospace'
  },
  Button: {
    textColorPrimary: '#1E1E1E',
    textColorHoverPrimary: '#1E1E1E',
    textColorPressedPrimary: '#1E1E1E',
    textColorFocusPrimary: '#1E1E1E',
    colorPrimary: museColors.dark.primary,
    colorHoverPrimary: museColors.dark.primaryHover,
    colorPressedPrimary: museColors.dark.primaryPressed,
    colorFocusPrimary: museColors.dark.primary,
    borderPrimary: '1px solid transparent',
    borderHoverPrimary: '1px solid transparent',
    borderPressedPrimary: '1px solid transparent',
    borderFocusPrimary: '1px solid transparent',
    rippleColorPrimary: museColors.dark.primary,

    textColorGhost: museColors.dark.primary,
    textColorGhostHover: museColors.dark.primaryHover,
    textColorGhostPressed: museColors.dark.primaryPressed,
    borderGhost: `1px solid ${museColors.dark.primary}`,
    borderGhostHover: `1px solid ${museColors.dark.primaryHover}`,
    borderGhostPressed: `1px solid ${museColors.dark.primaryPressed}`,

    borderRadiusMedium: '8px',
    borderRadiusSmall: '6px',
    borderRadiusLarge: '10px',

    fontWeightStrong: '500'
  },
  Card: {
    color: museColors.dark.bgCard,
    colorModal: museColors.dark.bgModal,
    colorPopover: museColors.dark.bgCard,
    colorTarget: museColors.dark.bgCard,
    colorEmbedded: museColors.dark.bgSecondary,
    borderColor: museColors.dark.border,
    borderRadius: '12px',
    boxShadow: '0 1px 3px rgba(0, 0, 0, 0.3)',
    titleFontSizeMedium: '16px',
    titleFontWeight: '500'
  },
  Input: {
    color: museColors.dark.bgTertiary,
    colorFocus: museColors.dark.bgTertiary,
    colorDisabled: museColors.dark.bgSecondary,
    border: `1px solid ${museColors.dark.border}`,
    borderHover: `1px solid ${museColors.dark.borderLight}`,
    borderFocus: `1px solid ${museColors.dark.primary}`,
    borderDisabled: `1px solid ${museColors.dark.border}`,
    boxShadowFocus: `0 0 0 2px rgba(77, 168, 255, 0.3)`,
    caretColor: museColors.dark.primary,
    borderRadius: '8px',
    heightMedium: '40px'
  },
  Select: {
    peers: {
      InternalSelection: {
        color: museColors.dark.bgTertiary,
        colorActive: museColors.dark.bgTertiary,
        colorDisabled: museColors.dark.bgSecondary,
        border: `1px solid ${museColors.dark.border}`,
        borderHover: `1px solid ${museColors.dark.borderLight}`,
        borderActive: `1px solid ${museColors.dark.primary}`,
        borderFocus: `1px solid ${museColors.dark.primary}`,
        boxShadowFocus: `0 0 0 2px rgba(77, 168, 255, 0.3)`,
        boxShadowActive: `0 0 0 2px rgba(77, 168, 255, 0.3)`,
        borderRadius: '8px',
        textColor: museColors.dark.textPrimary,
        placeholderColor: museColors.dark.textTertiary,
        arrowColor: museColors.dark.textSecondary
      },
      InternalSelectMenu: {
        color: museColors.dark.bgCard,
    optionColorPending: 'rgba(77, 168, 255, 0.06)',
        optionColorActive: 'rgba(77, 168, 255, 0.15)',
        optionTextColorActive: museColors.dark.primary,
        optionCheckColor: museColors.dark.primary,
        borderRadius: '8px'
      }
    }
  },
  Menu: {
    color: 'transparent',
    itemColorHover: 'rgba(77, 168, 255, 0.06)',
    itemColorActive: 'rgba(77, 168, 255, 0.12)',
    itemColorActiveHover: 'rgba(77, 168, 255, 0.16)',
    itemColorActiveCollapsed: 'rgba(77, 168, 255, 0.12)',
    itemTextColor: museColors.dark.textSecondary,
    itemTextColorHover: museColors.dark.textPrimary,
    itemTextColorActive: museColors.dark.primary,
    itemTextColorActiveHover: museColors.dark.primary,
    itemTextColorChildActive: museColors.dark.primary,
    itemIconColor: museColors.dark.textSecondary,
    itemIconColorHover: museColors.dark.textPrimary,
    itemIconColorActive: museColors.dark.primary,
    itemIconColorActiveHover: museColors.dark.primary,
    itemIconColorCollapsed: museColors.dark.textSecondary,
    arrowColor: museColors.dark.textSecondary,
    arrowColorHover: museColors.dark.textPrimary,
    arrowColorActive: museColors.dark.primary,
    arrowColorChildActive: museColors.dark.primary,
    borderRadius: '8px',
    itemHeight: '40px'
  },
  Modal: {
    color: museColors.dark.bgModal,
    textColor: museColors.dark.textPrimary,
    borderRadius: '16px',
    boxShadow: '0 12px 40px rgba(0, 0, 0, 0.5)'
  },
  Message: {
    borderRadius: '8px',
    boxShadow: '0 4px 16px rgba(0, 0, 0, 0.4)'
  },
  Tag: {
    color: museColors.dark.bgTertiary,
    textColor: museColors.dark.textSecondary,
    border: `1px solid ${museColors.dark.border}`,
    borderRadius: '6px',
    colorPrimary: 'rgba(77, 168, 255, 0.15)',
    textColorPrimary: museColors.dark.primary,
    borderPrimary: `1px solid rgba(77, 168, 255, 0.3)`,
    closeIconColor: museColors.dark.textTertiary,
    closeIconColorHover: museColors.dark.textPrimary
  },
  Avatar: {
    color: museColors.dark.bgTertiary,
    borderRadius: '50%'
  },
  Tabs: {
    tabTextColorLine: museColors.dark.textSecondary,
    tabTextColorActiveLine: museColors.dark.primary,
    tabTextColorHoverLine: museColors.dark.textPrimary,
    barColor: museColors.dark.primary,
    tabBorderColor: museColors.dark.border,
    tabGapMediumLine: '24px',
    tabPaddingMediumLine: '12px 0',
    tabFontWeightActive: '500'
  },
  Scrollbar: {
    color: 'rgba(77, 168, 255, 0.12)',
    colorHover: 'rgba(77, 168, 255, 0.24)',
    width: '8px',
    borderRadius: '4px'
  },
  Drawer: {
    color: museColors.dark.bgSecondary,
    textColor: museColors.dark.textPrimary,
    borderRadius: '0',
    boxShadow: '-4px 0 24px rgba(0, 0, 0, 0.4)'
  },
  Switch: {
    railColor: museColors.dark.border,
    railColorActive: museColors.dark.primary,
    buttonColor: '#E3E3E3',
    boxShadowFocus: `0 0 0 2px rgba(77, 168, 255, 0.3)`
  },
  Slider: {
    fillColor: museColors.dark.primary,
    fillColorHover: museColors.dark.primaryHover,
    railColor: museColors.dark.border,
    handleColor: museColors.dark.primary,
    dotBorderActive: `2px solid ${museColors.dark.primary}`
  },
  Progress: {
    fillColor: museColors.dark.primary,
    railColor: museColors.dark.border,
    textColorCircle: museColors.dark.textPrimary
  },
  List: {
    color: 'transparent',
    colorHover: 'rgba(255, 255, 255, 0.04)',
    borderColor: museColors.dark.border,
    textColor: museColors.dark.textPrimary
  },
  Empty: {
    textColor: museColors.dark.textTertiary,
    iconColor: museColors.dark.textTertiary
  },
  Tooltip: {
    color: museColors.dark.bgCard,
    textColor: museColors.dark.textPrimary,
    borderRadius: '8px',
    boxShadow: '0 4px 16px rgba(0, 0, 0, 0.4)'
  },
  Popover: {
    color: museColors.dark.bgCard,
    textColor: museColors.dark.textPrimary,
    borderRadius: '12px',
    boxShadow: '0 8px 24px rgba(0, 0, 0, 0.4)'
  },
  Dropdown: {
    color: museColors.dark.bgCard,
    optionColorHover: 'rgba(77, 168, 255, 0.06)',
    optionColorActive: 'rgba(77, 168, 255, 0.12)',
    optionTextColor: museColors.dark.textSecondary,
    optionTextColorHover: museColors.dark.textPrimary,
    optionTextColorActive: museColors.dark.primary,
    borderRadius: '8px',
    boxShadow: '0 8px 24px rgba(0, 0, 0, 0.4)'
  }
};

// 明亮主题覆盖
export const lightThemeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: museColors.light.primary,
    primaryColorHover: museColors.light.primaryHover,
    primaryColorPressed: museColors.light.primaryPressed,
    primaryColorSuppl: museColors.light.primarySuppl,

    infoColor: museColors.light.info,
    infoColorHover: '#4285F4',
    infoColorPressed: '#1565C0',

    successColor: museColors.light.success,
    successColorHover: '#34A853',
    successColorPressed: '#137333',

    warningColor: museColors.light.warning,
    warningColorHover: '#FBBC04',
    warningColorPressed: '#E37400',

    errorColor: museColors.light.error,
    errorColorHover: '#EA4335',
    errorColorPressed: '#B31412',

    textColorBase: museColors.light.textPrimary,
    textColor1: museColors.light.textPrimary,
    textColor2: museColors.light.textSecondary,
    textColor3: museColors.light.textTertiary,
    textColorDisabled: museColors.light.textDisabled,

    placeholderColor: museColors.light.textTertiary,
    placeholderColorDisabled: museColors.light.textDisabled,

    iconColor: museColors.light.textSecondary,
    iconColorHover: museColors.light.textPrimary,
    iconColorPressed: museColors.light.primary,
    iconColorDisabled: museColors.light.textDisabled,

    dividerColor: museColors.light.border,
    borderColor: museColors.light.border,

    closeIconColor: museColors.light.textSecondary,
    closeIconColorHover: museColors.light.textPrimary,
    closeIconColorPressed: museColors.light.primary,

    clearColor: museColors.light.textTertiary,
    clearColorHover: museColors.light.textSecondary,
    clearColorPressed: museColors.light.textPrimary,

    scrollbarColor: 'rgba(0, 0, 0, 0.12)',
    scrollbarColorHover: 'rgba(0, 0, 0, 0.2)',

    progressRailColor: museColors.light.border,

    bodyColor: museColors.light.bgPrimary,
    cardColor: museColors.light.bgCard,
    modalColor: museColors.light.bgModal,
    popoverColor: museColors.light.bgCard,
    tableColor: museColors.light.bgCard,
    inputColor: museColors.light.bgSecondary,
    inputColorDisabled: museColors.light.bgTertiary,
    actionColor: museColors.light.bgTertiary,
    tagColor: museColors.light.bgTertiary,
    avatarColor: museColors.light.bgTertiary,
    invertedColor: museColors.light.bgCard,

    hoverColor: 'rgba(0, 0, 0, 0.04)',
    pressedColor: 'rgba(0, 0, 0, 0.08)',

    boxShadow1: '0 1px 3px rgba(0, 0, 0, 0.08)',
    boxShadow2: '0 4px 12px rgba(0, 0, 0, 0.1)',
    boxShadow3: '0 8px 24px rgba(0, 0, 0, 0.12)',

    borderRadius: '8px',
    borderRadiusSmall: '4px',

    fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Noto Sans SC", sans-serif',
    fontFamilyMono: '"JetBrains Mono", "Fira Code", "Source Code Pro", monospace'
  },
  Button: {
    textColorPrimary: '#FFFFFF',
    textColorHoverPrimary: '#FFFFFF',
    textColorPressedPrimary: '#FFFFFF',
    textColorFocusPrimary: '#FFFFFF',
    colorPrimary: museColors.light.primary,
    colorHoverPrimary: museColors.light.primaryHover,
    colorPressedPrimary: museColors.light.primaryPressed,
    colorFocusPrimary: museColors.light.primary,
    borderPrimary: '1px solid transparent',
    borderHoverPrimary: '1px solid transparent',
    borderPressedPrimary: '1px solid transparent',
    borderFocusPrimary: '1px solid transparent',
    rippleColorPrimary: museColors.light.primary,

    textColorGhost: museColors.light.primary,
    textColorGhostHover: museColors.light.primaryHover,
    textColorGhostPressed: museColors.light.primaryPressed,
    borderGhost: `1px solid ${museColors.light.primary}`,
    borderGhostHover: `1px solid ${museColors.light.primaryHover}`,
    borderGhostPressed: `1px solid ${museColors.light.primaryPressed}`,

    borderRadiusMedium: '8px',
    borderRadiusSmall: '6px',
    borderRadiusLarge: '10px',

    fontWeightStrong: '500'
  },
  Card: {
    color: museColors.light.bgCard,
    colorModal: museColors.light.bgModal,
    colorPopover: museColors.light.bgCard,
    colorTarget: museColors.light.bgCard,
    colorEmbedded: museColors.light.bgSecondary,
    borderColor: museColors.light.border,
    borderRadius: '12px',
    boxShadow: '0 1px 3px rgba(0, 0, 0, 0.08)',
    titleFontSizeMedium: '16px',
    titleFontWeight: '500'
  },
  Input: {
    color: museColors.light.bgSecondary,
    colorFocus: museColors.light.bgSecondary,
    colorDisabled: museColors.light.bgTertiary,
    border: `1px solid ${museColors.light.border}`,
    borderHover: `1px solid ${museColors.light.borderLight}`,
    borderFocus: `1px solid ${museColors.light.primary}`,
    borderDisabled: `1px solid ${museColors.light.border}`,
    boxShadowFocus: `0 0 0 2px rgba(37, 99, 235, 0.2)`,
    caretColor: museColors.light.primary,
    borderRadius: '8px',
    heightMedium: '40px'
  },
  Select: {
    peers: {
      InternalSelection: {
        color: museColors.light.bgSecondary,
        colorActive: museColors.light.bgSecondary,
        colorDisabled: museColors.light.bgTertiary,
        border: `1px solid ${museColors.light.border}`,
        borderHover: `1px solid ${museColors.light.borderLight}`,
        borderActive: `1px solid ${museColors.light.primary}`,
        borderFocus: `1px solid ${museColors.light.primary}`,
        boxShadowFocus: `0 0 0 2px rgba(37, 99, 235, 0.2)`,
        boxShadowActive: `0 0 0 2px rgba(37, 99, 235, 0.2)`,
        borderRadius: '8px',
        textColor: museColors.light.textPrimary,
        placeholderColor: museColors.light.textTertiary,
        arrowColor: museColors.light.textSecondary
      },
      InternalSelectMenu: {
        color: museColors.light.bgCard,
        optionColorPending: 'rgba(0, 0, 0, 0.04)',
        optionColorActive: 'rgba(37, 99, 235, 0.1)',
        optionTextColorActive: museColors.light.primary,
        optionCheckColor: museColors.light.primary,
        borderRadius: '8px'
      }
    }
  },
  Menu: {
    color: 'transparent',
    itemColorHover: 'rgba(0, 0, 0, 0.04)',
    itemColorActive: 'rgba(37, 99, 235, 0.08)',
    itemColorActiveHover: 'rgba(37, 99, 235, 0.12)',
    itemColorActiveCollapsed: 'rgba(37, 99, 235, 0.08)',
    itemTextColor: museColors.light.textSecondary,
    itemTextColorHover: museColors.light.textPrimary,
    itemTextColorActive: museColors.light.primary,
    itemTextColorActiveHover: museColors.light.primary,
    itemTextColorChildActive: museColors.light.primary,
    itemIconColor: museColors.light.textSecondary,
    itemIconColorHover: museColors.light.textPrimary,
    itemIconColorActive: museColors.light.primary,
    itemIconColorActiveHover: museColors.light.primary,
    itemIconColorCollapsed: museColors.light.textSecondary,
    arrowColor: museColors.light.textSecondary,
    arrowColorHover: museColors.light.textPrimary,
    arrowColorActive: museColors.light.primary,
    arrowColorChildActive: museColors.light.primary,
    borderRadius: '8px',
    itemHeight: '40px'
  },
  Modal: {
    color: museColors.light.bgModal,
    textColor: museColors.light.textPrimary,
    borderRadius: '16px',
    boxShadow: '0 12px 40px rgba(0, 0, 0, 0.15)'
  },
  Message: {
    borderRadius: '8px',
    boxShadow: '0 4px 16px rgba(0, 0, 0, 0.1)'
  },
  Tag: {
    color: museColors.light.bgTertiary,
    textColor: museColors.light.textSecondary,
    border: `1px solid ${museColors.light.border}`,
    borderRadius: '6px',
    colorPrimary: 'rgba(37, 99, 235, 0.08)',
    textColorPrimary: museColors.light.primary,
    borderPrimary: `1px solid rgba(37, 99, 235, 0.3)`,
    closeIconColor: museColors.light.textTertiary,
    closeIconColorHover: museColors.light.textPrimary
  },
  Avatar: {
    color: museColors.light.bgTertiary,
    borderRadius: '50%'
  },
  Tabs: {
    tabTextColorLine: museColors.light.textSecondary,
    tabTextColorActiveLine: museColors.light.primary,
    tabTextColorHoverLine: museColors.light.textPrimary,
    barColor: museColors.light.primary,
    tabBorderColor: museColors.light.border,
    tabGapMediumLine: '24px',
    tabPaddingMediumLine: '12px 0',
    tabFontWeightActive: '500'
  },
  Scrollbar: {
    color: 'rgba(0, 0, 0, 0.12)',
    colorHover: 'rgba(0, 0, 0, 0.2)',
    width: '8px',
    borderRadius: '4px'
  },
  Drawer: {
    color: museColors.light.bgSecondary,
    textColor: museColors.light.textPrimary,
    borderRadius: '0',
    boxShadow: '-4px 0 24px rgba(0, 0, 0, 0.1)'
  },
  Switch: {
    railColor: museColors.light.border,
    railColorActive: museColors.light.primary,
    buttonColor: '#FFFFFF',
    boxShadowFocus: `0 0 0 2px rgba(37, 99, 235, 0.2)`
  },
  Slider: {
    fillColor: museColors.light.primary,
    fillColorHover: museColors.light.primaryHover,
    railColor: museColors.light.border,
    handleColor: museColors.light.primary,
    dotBorderActive: `2px solid ${museColors.light.primary}`
  },
  Progress: {
    fillColor: museColors.light.primary,
    railColor: museColors.light.border,
    textColorCircle: museColors.light.textPrimary
  },
  List: {
    color: 'transparent',
    colorHover: 'rgba(0, 0, 0, 0.03)',
    borderColor: museColors.light.border,
    textColor: museColors.light.textPrimary
  },
  Empty: {
    textColor: museColors.light.textTertiary,
    iconColor: museColors.light.textTertiary
  },
  Tooltip: {
    color: museColors.light.bgCard,
    textColor: museColors.light.textPrimary,
    borderRadius: '8px',
    boxShadow: '0 4px 16px rgba(0, 0, 0, 0.1)'
  },
  Popover: {
    color: museColors.light.bgCard,
    textColor: museColors.light.textPrimary,
    borderRadius: '12px',
    boxShadow: '0 8px 24px rgba(0, 0, 0, 0.1)'
  },
  Dropdown: {
    color: museColors.light.bgCard,
    optionColorHover: 'rgba(0, 0, 0, 0.04)',
    optionColorActive: 'rgba(37, 99, 235, 0.08)',
    optionTextColor: museColors.light.textSecondary,
    optionTextColorHover: museColors.light.textPrimary,
    optionTextColorActive: museColors.light.primary,
    borderRadius: '8px',
    boxShadow: '0 8px 24px rgba(0, 0, 0, 0.1)'
  }
};
