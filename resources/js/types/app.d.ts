export interface User {
    id: number
    name: string
    email: string
    role: string
}

export interface SharedData {
    pageTitle: string;
    auth: {
        user: User | null;
    };
    // Add other specific props for this page if any
}