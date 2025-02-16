declare global {
    interface Window {
        __DATA__: Record<string, unknown>;
    }
}

export {};
