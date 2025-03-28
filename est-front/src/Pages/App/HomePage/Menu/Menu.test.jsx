import React, { act } from 'react';
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { message } from 'antd';
import { MemoryRouter, useNavigate } from 'react-router-dom';
import AppMenu from './Menu';

jest.mock('antd', () => {
    const antd = jest.requireActual('antd');
    return {
        ...antd,
        message: {
            useMessage: jest.fn(),
        },
    };
});

jest.mock('est_proxy_api', () => ({
    UserApi: jest.fn().mockImplementation(() => ({
        getSelf: jest.fn(),
        logout: jest.fn(),
    })),
}));

jest.mock('react-router-dom', () => ({
    ...jest.requireActual('react-router-dom'),
    useNavigate: jest.fn(),
}));

window.matchMedia = window.matchMedia || function() {
    return {
        matches: false,
        addListener: function() {},
        removeListener: function() {}
    };
};

describe('AppMenu Component', () => {
    let mockGetSelf;
    let mockLogout;
    let mockMessageApi;
    let mockNavigate;

    beforeEach(() => {
        mockMessageApi = {
            open: jest.fn(),
        };
        message.useMessage.mockReturnValue([mockMessageApi, jest.fn()]);

        mockNavigate = jest.fn();
        useNavigate.mockReturnValue(mockNavigate);

        mockGetSelf = jest.fn();
        mockLogout = jest.fn();

        const mockUserApi = {
            apiClient: {
                basePath: '',
                defaultHeaders: {},
            },
            getSelf: mockGetSelf,
            logout: mockLogout,
        };
        require('est_proxy_api').UserApi.mockImplementation(() => mockUserApi);
    });

    afterEach(() => {
        jest.clearAllMocks();
    });

    it('renders nothing when user data is not loaded', async () => {
        mockGetSelf.mockImplementation(() => new Promise(() => {}));

        await act(async () => {
            render(
                <MemoryRouter>
                    <AppMenu />
                </MemoryRouter>
            );
        });

        await waitFor(() => {
            expect(screen.getByText("Loading...")).toBeInTheDocument();
        });
    });

    it('loads user data and renders menu', async () => {
        const mockUserData = { username: 'testuser' };
        mockGetSelf.mockResolvedValue(mockUserData);

        await act(async () => {
            render(
                <MemoryRouter>
                    <AppMenu />
                </MemoryRouter>
            );
        });

        await waitFor(() => {
            expect(mockGetSelf).toHaveBeenCalled();

            expect(screen.getByText('e-Sketch')).toBeInTheDocument();
            expect(screen.getByText('testuser')).toBeInTheDocument();
            expect(screen.getByText('Настройки')).toBeInTheDocument();
            expect(screen.getByText('Доски')).toBeInTheDocument();
            expect(screen.getByText('Мои доски')).toBeInTheDocument();
            expect(screen.getByText('Поделились со мной')).toBeInTheDocument();
            expect(screen.getByText('Выйти')).toBeInTheDocument();
        });
    });

    it('shows error message when user data loading fails', async () => {
        const error = new Error('Failed to load user data');
        mockGetSelf.mockRejectedValue(error);

        await act(async () => {
            render(
                <MemoryRouter>
                    <AppMenu />
                </MemoryRouter>
            );
        });

        await waitFor(() => {
            expect(mockMessageApi.open).toHaveBeenCalledWith({
                type: 'error',
                content: error,
            });
        });
    });

    it('navigates to home page on successful logout', async () => {
        const mockUserData = { username: 'testuser' };
        mockGetSelf.mockResolvedValue(mockUserData);
        mockLogout.mockResolvedValue({});

        await act(async () => {
            render(
                <MemoryRouter>
                    <AppMenu />
                </MemoryRouter>
            );
        });

        await waitFor(() => {
            expect(screen.getByText('Выйти')).toBeInTheDocument();
        });

        fireEvent.click(screen.getByText('Выйти'));

        await waitFor(() => {
            expect(mockLogout).toHaveBeenCalled();
            expect(mockNavigate).toHaveBeenCalledWith('/');
        });
    });

    it('shows error message when logout fails', async () => {
        const mockUserData = { username: 'testuser' };
        mockGetSelf.mockResolvedValue(mockUserData);
        const error = new Error('Logout failed');
        mockLogout.mockRejectedValue(error);

        await act(async () => {
            render(
                <MemoryRouter>
                    <AppMenu />
                </MemoryRouter>
            );
        });

        await waitFor(() => {
            expect(screen.getByText('Выйти')).toBeInTheDocument();
        });

        fireEvent.click(screen.getByText('Выйти'));

        await waitFor(() => {
            expect(mockMessageApi.open).toHaveBeenCalledWith({
                type: 'error',
                content: error,
            });
        });
    });

    it('renders correct menu items with icons', async () => {
        const mockUserData = { username: 'testuser' };
        mockGetSelf.mockResolvedValue(mockUserData);

        await act(async () => {
            render(
                <MemoryRouter>
                    <AppMenu />
                </MemoryRouter>
            );
        });
        
        await waitFor(() => {
            expect(screen.getByText('testuser')).toBeInTheDocument();

            expect(screen.getByLabelText('user')).toBeInTheDocument();
            expect(screen.getByLabelText('setting')).toBeInTheDocument();
            expect(screen.getByLabelText('appstore')).toBeInTheDocument();
            expect(screen.getByLabelText('link')).toBeInTheDocument();
            expect(screen.getByLabelText('logout')).toBeInTheDocument();
        });
    });

    it('applies correct layout and styling when rendered', async () => {
        const mockUserData = { username: 'testuser' };
        mockGetSelf.mockResolvedValue(mockUserData);

        await act(async () => {
        render(
            <MemoryRouter>
            <AppMenu />
            </MemoryRouter>
        );
        });

        await waitFor(() => {
            expect(screen.getByText('testuser')).toBeInTheDocument();

            const menu = screen.getByRole('menu');
            expect(menu).toBeInTheDocument();
            expect(menu).toHaveStyle({
                width: '256px',
                minHeight: '600px'
            });

            const flexContainer = screen.getByRole('button', { name: /выйти/i }).parentElement;
            expect(flexContainer).toHaveStyle({
                width: '100%'
            });
        });
    });
});