export interface User {
    id: number
    name: string
    email: string
    role: string
    is_active?: boolean
    isSuperAdmin?: boolean
    isAdmin?: boolean
    permissions?: string[]
    roles?: Array<{
        id: number
        name: string
        slug?: string
    }>
}

export interface SharedData {
    pageTitle: string;
    appVersion: string;
    auth: {
        user: User | null;
    };
    // Index signature to satisfy PageProps constraint
    [key: string]: any;
}