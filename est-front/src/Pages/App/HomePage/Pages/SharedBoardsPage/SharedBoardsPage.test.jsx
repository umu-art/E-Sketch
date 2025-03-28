import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { message } from 'antd';
import { MemoryRouter, useNavigate } from 'react-router-dom';
import SharedBoardsPage from './SharedBoardsPage';

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

window.matchMedia = window.matchMedia || function() {
    return {
        matches: false,
        addListener: function() {},
        removeListener: function() {}
    };
};

describe('SharedBoardsPage Component', () => {
    let mockList;
    let mockMessageApi;
    let mockNavigate;

    const mockBoards = {
        shared: [
            { id: '1', name: 'Shared Board 1' },
            { id: '2', name: 'Shared Board 2' },
        ]
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
            <SharedBoardsPage />
        );

        expect(screen.getByRole('img', { name: /loading/i })).toBeInTheDocument();
    });

    it('loads and displays shared boards correctly', async () => {
        mockList.mockResolvedValue(mockBoards);

        render(
            <MemoryRouter>
                <SharedBoardsPage />
            </MemoryRouter>
        );

        await waitFor(() => {
            expect(screen.getByText('Поделились со мной')).toBeInTheDocument();
            expect(screen.getByText('Shared Board 1')).toBeInTheDocument();
            expect(screen.getByText('Shared Board 2')).toBeInTheDocument();
        });
    });

    it('shows no boards message when shared boards are empty', async () => {
        mockList.mockResolvedValue({ shared: [] });

        render(
            <MemoryRouter>
                <SharedBoardsPage />
            </MemoryRouter>
        );

        await waitFor(() => {
            expect(screen.getByText('Поделились со мной')).toBeInTheDocument();
            expect(screen.queryByTestId(/board-card-/)).not.toBeInTheDocument();
        });
    });

    it('handles unauthorized error by redirecting to signin', async () => {
        const error = { statusCode: 401, rawResponse: 'Unauthorized' };
        mockList.mockRejectedValue(error);

        render(
            <MemoryRouter>
                <SharedBoardsPage />
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

    it('applies correct styling and layout', async () => {
        mockList.mockResolvedValue(mockBoards);

        render(
            <MemoryRouter>
                <SharedBoardsPage />
            </MemoryRouter>
        );

        await waitFor(() => {
            const container = screen.getByText('Поделились со мной').parentElement.parentElement;
            expect(container).toHaveStyle({
                padding: '20px 50px',
            });

            const boardsContainer = screen.getByText('Shared Board 1').parentElement;
            expect(boardsContainer).toHaveStyle({
                width: 'auto',
            });
        });
    });
});