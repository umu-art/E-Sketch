import React from 'react';
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { message } from 'antd';
import { MemoryRouter, useNavigate } from 'react-router-dom';
import MyBoardsPage from './MyBoardsPage';


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
    BoardApi: jest.fn().mockImplementation(() => ({
        list: jest.fn(),
        apiClient: {
            basePath: '',
            defaultHeaders: {},
        },
    })),
}));

jest.mock('react-router-dom', () => ({
    ...jest.requireActual('react-router-dom'),
    useNavigate: jest.fn(),
}));

jest.mock("../../BoardCard/BoardCard", () => ({ data }) => (
    <div data-testid={`board-card-${data.id}`}>{data.name}</div>
));

jest.mock('../../CreateBoardForm/CreateBoardForm', () => () => (
    <div data-testid="create-board-form">Create Board Form</div>
));

window.matchMedia = window.matchMedia || function() {
    return {
        matches: false,
        addListener: function() {},
        removeListener: function() {}
    };
};

describe('MyBoardsPage Component', () => {
    let mockList;
    let mockMessageApi;
    let mockNavigate;

    const mockBoards = {
        mine: [
            { id: '1', name: 'My Board 1' },
            { id: '2', name: 'My Board 2' },
        ],
        recent: [
            { id: '3', name: 'Recent Board 1' },
            { id: '4', name: 'Recent Board 2' },
        ],
    };

    beforeEach(() => {
        mockMessageApi = {
            open: jest.fn(),
        };
        message.useMessage.mockReturnValue([mockMessageApi, jest.fn()]);

        mockNavigate = jest.fn();
        useNavigate.mockReturnValue(mockNavigate);

        mockList = jest.fn();
        require('est_proxy_api').BoardApi.mockImplementation(() => ({
            apiClient: {
                basePath: '',
                defaultHeaders: {},
            },
            list: mockList,
        }));
    });

    afterEach(() => {
        jest.clearAllMocks();
    });

    it('renders loading state initially', async () => {
        mockList.mockImplementation(() => new Promise(() => {}));

        render(
            <MyBoardsPage />
        );

        expect(screen.getByRole('img', { name: /loading/i })).toBeInTheDocument();
    });

    it('loads and displays boards correctly', async () => {
        mockList.mockResolvedValue(mockBoards);

        render(
            <MemoryRouter>
                <MyBoardsPage />
            </MemoryRouter>
        );

        await waitFor(() => {
            expect(screen.getByText('Мои доски')).toBeInTheDocument();
            expect(screen.getByText('My Board 1')).toBeInTheDocument();
            expect(screen.getByText('My Board 2')).toBeInTheDocument();
        });
    });

    it('shows only my boards when no recent boards', async () => {
        mockList.mockResolvedValue({
            mine: mockBoards.mine,
            recent: [],
        });

        render(
            <MyBoardsPage />
        );

        await waitFor(() => {
            expect(screen.getByText('Мои доски')).toBeInTheDocument();
            expect(screen.queryByText('Недавние')).not.toBeInTheDocument();
        });
    });

    it('handles unauthorized error by redirecting to signin', async () => {
        const error = { statusCode: 401, rawResponse: 'Unauthorized' };
        mockList.mockRejectedValue(error);

        render(
            <MemoryRouter>
                <MyBoardsPage />
            </MemoryRouter>
        );

        await waitFor(() => {
            expect(mockNavigate).toHaveBeenCalledWith('/auth/signin');
            expect(mockMessageApi.open).toHaveBeenCalledWith({
                type: 'error',
                content: error.rawResponse,
            });
        });
    });

    it('opens create board modal when button clicked', async () => {
        mockList.mockResolvedValue({ mine: [], recent: [] });

        render(
            <MemoryRouter>
                <MyBoardsPage />
            </MemoryRouter>
        );

        fireEvent.click(screen.getByText('Новая доска'));

        await waitFor(() => {
            expect(screen.getByText('Создать доску')).toBeInTheDocument();
            expect(screen.getByTestId('create-board-form')).toBeInTheDocument();
        });
    });

    it('closes create board modal when cancelled', async () => {
        mockList.mockResolvedValue({ mine: [], recent: [] });

        render(
            <MemoryRouter>
                <MyBoardsPage />
            </MemoryRouter>
        );

        fireEvent.click(screen.getByText('Новая доска'));
        fireEvent.click(screen.getByRole('button', { name: /close/i }));

        

        await waitFor(() => {
            expect(screen.queryByText('Создать доску')).not.toBeVisible();
        });
    });

    it('applies correct styling and layout', async () => {
        mockList.mockResolvedValue(mockBoards);

        render(
            <MemoryRouter>
                <MyBoardsPage />
            </MemoryRouter>
        );

        await waitFor(() => {
            const container = screen.getByText('Мои доски').parentElement.parentElement;
            expect(container).toHaveStyle({
                padding: '20px 50px',
            });

            const boardsContainer = screen.getByText('My Board 1').parentElement;
            expect(boardsContainer).toHaveStyle({
                width: 'auto',
            });
        });
    });

    it('does not show recent section when recent boards are empty', async () => {
        mockList.mockResolvedValue({
            mine: mockBoards.mine,
            recent: [],
        });

        render(
            <MemoryRouter>
                <MyBoardsPage />
            </MemoryRouter>
        );

        await waitFor(() => {
            expect(screen.queryByText('Недавние')).not.toBeInTheDocument();
        });
    });

    it('does not show recent section when my boards are less than 4', async () => {
        mockList.mockResolvedValue({
            mine: [{ id: '1', name: 'Single Board' }],
            recent: mockBoards.recent,
        });

        render(
            <MemoryRouter>
                <MyBoardsPage />
            </MemoryRouter>
        );

        await waitFor(() => {
            expect(screen.queryByText('Недавние')).not.toBeInTheDocument();
        });
    });
});