import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import UpdateBoardForm from './UpdateBoardForm';
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
        update: jest.fn(),
    })),
}));

describe('UpdateBoardForm Component', () => {
    const mockData = {
        id: '1',
        name: 'Test Board',
        description: 'Test Description'
    };
    const mockOnDataChange = jest.fn();
    const mockCloseModal = jest.fn();
    
    let mockUpdate;
    let mockMessageApi;

    beforeEach(() => {
        mockMessageApi = {
            open: jest.fn(),
        };
        message.useMessage.mockReturnValue([mockMessageApi, jest.fn()]);

        mockUpdate = jest.fn();
        const mockBoardApi = {
            apiClient: {
                basePath: '',
                defaultHeaders: {},
            },
            update: mockUpdate,
        };
        require('est_proxy_api').BoardApi.mockImplementation(() => mockBoardApi);
    });

    afterEach(() => {
        jest.clearAllMocks();
    });

    it('renders form with initial data and submit button', () => {
        render(
            <UpdateBoardForm 
                data={mockData} 
                onDataChange={mockOnDataChange} 
                closeModal={mockCloseModal} 
            />
        );
        
        expect(screen.getByLabelText('Название')).toBeInTheDocument();
        expect(screen.getByDisplayValue('Test Board')).toBeInTheDocument();
        expect(screen.getByLabelText('Описание')).toBeInTheDocument();
        expect(screen.getAllByText('Test Description')[0]).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Сохранить' })).toBeInTheDocument();
    });

    it('shows error when name is empty', async () => {
        render(
            <UpdateBoardForm 
                data={mockData} 
                onDataChange={mockOnDataChange} 
                closeModal={mockCloseModal} 
            />
        );
        
        fireEvent.change(screen.getByPlaceholderText('Введите название доски...'), {
            target: { value: '' }
        });
        fireEvent.click(screen.getByRole('button', { name: 'Сохранить' }));
        
        await waitFor(() => {
            expect(screen.getByText('Пожалуйста, введите название доски!')).toBeInTheDocument();
        });
    });

    it('submits form with updated data', async () => {
        const newData = { ...mockData, name: 'Updated Board' };
        mockUpdate.mockResolvedValueOnce(newData);
        
        render(
            <UpdateBoardForm 
                data={mockData} 
                onDataChange={mockOnDataChange} 
                closeModal={mockCloseModal} 
            />
        );
        
        fireEvent.change(screen.getByPlaceholderText('Введите название доски...'), {
            target: { value: 'Updated Board' }
        });
        fireEvent.click(screen.getByRole('button', { name: 'Сохранить' }));
        
        await waitFor(() => {
            expect(mockUpdate).toHaveBeenCalledWith('1', {
                createRequest: {
                    ...mockData,
                    name: 'Updated Board'
                }
            });
        });
    });

    it('calls onDataChange and closeModal on successful submission', async () => {
        const newData = { ...mockData, name: 'Updated Board' };
        mockUpdate.mockResolvedValueOnce(newData);
        
        render(
            <UpdateBoardForm 
                data={mockData} 
                onDataChange={mockOnDataChange} 
                closeModal={mockCloseModal} 
            />
        );
        
        fireEvent.change(screen.getByPlaceholderText('Введите название доски...'), {
            target: { value: 'Updated Board' }
        });
        fireEvent.click(screen.getByRole('button', { name: 'Сохранить' }));
        
        await waitFor(() => {
            expect(mockOnDataChange).toHaveBeenCalledWith(newData);
            expect(mockCloseModal).toHaveBeenCalled();
            expect(mockMessageApi.open).toHaveBeenCalledWith({
                type: 'success',
                content: 'Изменения сохранены!'
            });
        });
    });

    it('shows error message when API call fails', async () => {
        const error = {
            response: 'Error updating board'
        };
        mockUpdate.mockRejectedValueOnce(error);
        
        render(
            <UpdateBoardForm 
                data={mockData} 
                onDataChange={mockOnDataChange} 
                closeModal={mockCloseModal} 
            />
        );
        
        fireEvent.click(screen.getByRole('button', { name: 'Сохранить' }));
        
        await waitFor(() => {
            expect(mockMessageApi.open).toHaveBeenCalledWith({
                type: 'error',
                content: 'Error updating board'
            });
        });
    });

    it('does not call onDataChange or closeModal when API call fails', async () => {
        mockUpdate.mockRejectedValueOnce(new Error('Network error'));
        
        render(
            <UpdateBoardForm 
                data={mockData} 
                onDataChange={mockOnDataChange} 
                closeModal={mockCloseModal} 
            />
        );
        
        fireEvent.click(screen.getByRole('button', { name: 'Сохранить' }));
        
        await waitFor(() => {
            expect(mockOnDataChange).not.toHaveBeenCalled();
            expect(mockCloseModal).not.toHaveBeenCalled();
        });
    });
});