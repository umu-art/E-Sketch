import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import CreateBoardForm from './CreateBoardForm';
import { message } from 'antd';

window.matchMedia = window.matchMedia || function() {
    return {
        matches: false,
        addListener: function() {},
        removeListener: function() {}
    };
};

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
        apiClient: {
            basePath: '',
            defaultHeaders: {},
        },
        create: jest.fn(),
    })),
}));

beforeAll(() => {
    Object.defineProperty(window, 'location', {
        value: {
        reload: jest.fn()
        },
        writable: true
    });
});

describe('CreateBoardForm Component', () => {
    let mockCreate;
    let mockMessageApi;

    beforeEach(() => {
        mockMessageApi = {
            open: jest.fn(),
        };
        message.useMessage.mockReturnValue([mockMessageApi, jest.fn()]);

        mockCreate = jest.fn();
        const mockBoardApi = {
            apiClient: {
                basePath: '',
                defaultHeaders: {},
            },
            create: mockCreate,
        };
        require('est_proxy_api').BoardApi.mockImplementation(() => mockBoardApi);
    });

    afterEach(() => {
        jest.clearAllMocks();
    });

    it('renders form with all fields and submit button', () => {
        render(<CreateBoardForm />);
        
        expect(screen.getByLabelText('Название')).toBeInTheDocument();
        expect(screen.getByPlaceholderText('Введите название доски...')).toBeInTheDocument();
        expect(screen.getByLabelText('Описание')).toBeInTheDocument();
        expect(screen.getByPlaceholderText('Введите описание доски...')).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Создать' })).toBeInTheDocument();
    });

    it('shows error when name is empty', async () => {
        render(<CreateBoardForm />);
        
        fireEvent.click(screen.getByRole('button', { name: 'Создать' }));
        
        await waitFor(() => {
            expect(screen.getByText('Пожалуйста, введите название доски!')).toBeInTheDocument();
        });
    });

    it('submits form with valid data', async () => {
        mockCreate.mockResolvedValueOnce({});
        
        render(<CreateBoardForm />);
        
        fireEvent.change(screen.getByPlaceholderText('Введите название доски...'), {
            target: { value: 'Test Board' }
        });
        fireEvent.change(screen.getByPlaceholderText('Введите описание доски...'), {
            target: { value: 'Test Description' }
        });
        fireEvent.click(screen.getByRole('button'));
        
        await waitFor(() => {
            expect(mockCreate).toHaveBeenCalledWith({
                createRequest: {
                    name: 'Test Board',
                    description: 'Test Description',
                    linkSharedMode: 'none_by_link'
                }
            });
        });
    });

    it('shows success message and reloads page on successful submission', async () => {
        mockCreate.mockResolvedValueOnce({});
        
        render(<CreateBoardForm />);
        
        fireEvent.change(screen.getByPlaceholderText('Введите название доски...'), {
            target: { value: 'Test Board' }
        });
        fireEvent.click(screen.getByRole('button', { name: 'Создать' }));
        
        await Promise.resolve();
        
        await waitFor(() => {
            expect(mockMessageApi.open).toHaveBeenCalledWith({
                type: 'success',
                content: 'Доска создана!'
            });
            expect(window.location.reload).toHaveBeenCalled();
        });
    });

    it('shows error message when API call fails', async () => {
        const error = {
        response: {
            text: 'Error creating board'
        }
        };
        mockCreate.mockRejectedValueOnce(error);
        
        render(<CreateBoardForm />);
        
        fireEvent.change(screen.getByPlaceholderText('Введите название доски...'), {
            target: { value: 'Test Board' }
        });
        fireEvent.click(screen.getByRole('button', { name: 'Создать' }));
        
        await Promise.resolve();
        
        await waitFor(() => {
            expect(mockMessageApi.open).toHaveBeenCalledWith({
                type: 'error',
                content: 'Error creating board'
            });
        });
    });

    it('shows unknown error message when API call fails', async () => {
        mockCreate.mockRejectedValue(new Error('Network error'));
        
        render(<CreateBoardForm />);
        
        fireEvent.change(screen.getByPlaceholderText('Введите название доски...'), {
            target: { value: 'Test Board' }
        });
        fireEvent.click(screen.getByRole('button', { name: 'Создать' }));
        
        await Promise.resolve();
        
        await waitFor(() => {
            expect(mockMessageApi.open).toHaveBeenCalledWith({
                type: 'error',
                content: 'Произошла ошибка, попробуйте позже'
            });
        });
    });
});